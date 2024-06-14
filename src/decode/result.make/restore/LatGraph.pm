#!/usr/bin/perl 

#XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
#X                                                                  X
#X LatGraph.pl, Copyright(C) Gailius Raðkinis, 2020                 X
#X                                                                  X
#XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX

use strict;
use warnings;
use Encode qw(decode encode);
use File::Basename;
use List::Util qw[min max];
#use Data::Dumper qw(Dumper);

my $qou = decode("utf-8", qq/\xe2\x80\x9e/); # „
my $qcu = decode("utf-8", qq/\xe2\x80\x9c/); # “

# 124 S150
# BESI <keturioliktudeðimtu> ???
# 155 S218
# 90 tûkst. -> 90000
# 97 S150
# 337 mln. -> 337000000

sub del_int {
   my ( $array_ref, $n ) = @_;

   for(my $i=0; $i<scalar @$array_ref; $i++) {
      if ($$array_ref[$i] == $n) {
         splice @$array_ref, $i, 1;
         last;
         }
      }
}
#-----------------------------

my $dirname = dirname(__FILE__);
require "$dirname/lat_map.pl";
require "$dirname/lt.pl";

package LatGraph;

sub new {
   my $class = shift;
   my $self = {
      _name => shift,
      _frameRate => shift,
      _startTime => shift,
      _endTime   => shift,
      _frameBased => shift,
      _n_tm => {}, # in frames
      _in   => {},
      _out  => {},
      _e    => [],
   };
#   print "_startTime= $self->{_startTime}\n";
   $self->{_n_tm}->{0} = 0;
   bless $self, $class;
   return $self;
}

sub frame2time {
   my ( $self, $frame ) = @_;
   
   return $self->{_startTime} + $frame * $self->{_frameRate};
}
#-----------------------------

sub node2time {
   my ( $self, $node ) = @_;

   return $self->frame2time($self->{_n_tm}->{$node});
}
#-----------------------------

sub edge2time {
   my ( $self, $edge_id ) = @_;

   return abs($self->{_n_tm}->{$self->{_e}->[$edge_id]->{n2}} - $self->{_n_tm}->{$self->{_e}->[$edge_id]->{n1}}) * $self->{_frameRate};
}

sub find_entry_node {
   my ( $self ) = @_;

   foreach my $n (keys %{$self->{_n_tm}}) {
      return($n) if (!defined $self->{_in}->{$n}); # no in edges
      }
   return(-1);
}
#-----------------------------

sub find_exit_node {
   my ( $self ) = @_;

   foreach my $n (keys %{$self->{_n_tm}}) {
      return($n) if (!defined $self->{_out}->{$n}); # no out edges
      }
   return(-1);
}
#-----------------------------

# empty lattice consists of <eps>'es
sub is_empty {
   my ( $self ) = @_;

   for(my $i=0; $i<scalar @{$self->{_e}}; $i++) {
      return 0 if ($self->{_e}->[$i]->{word_id} == -1); # added edge
      return 0 if (main::int2word($self->{_e}->[$i]->{word_id}) ne '<eps>');
      }

   return 1;
   }
#-----------------------------

# used to read L1 & L2 (frame-based)
sub add_edge {
   my ( $self, $n1, $n2, $word_id, $digest, $nb_frames ) = @_;
  
   #print STDERR $digest."\n";

   my %edge = ( n1 => $n1, n2 => $n2, word_id => $word_id, digest => $digest, stt => 0 );
   push @{$self->{_e}}, \%edge;
   my $edge_id = (scalar @{$self->{_e}}) - 1;

   $self->{_n_tm}->{$n2} = $self->{_n_tm}->{$n1} + $nb_frames; 

   $self->{_out}->{$n1} = () if (!defined $self->{_out}->{$n1});
   push @{$self->{_out}->{$n1}}, $edge_id;
   $self->{_in}->{$n2} = () if (!defined $self->{_in}->{$n2});
   push @{$self->{_in}->{$n2}}, $edge_id;
}
#-----------------------------

# used to read L3 (phone-based)
sub add_phones {
   my ( $self, $phone_ids ) = @_;

   my @C = split("_", $phone_ids);
   my $transc='';
   for (my $p = 0; $p <= $#C; $p++) {
      $transc .= ' ' if ($p > 0);                   # space-separated phoneme string 
      $transc = $transc.main::int2phone($C[$p], 1); # remove_pd_info
      }
#   $transc = '<'.$transc.'>';
   my %edge = ( ph => $transc );
   push @{$self->{_e}}, \%edge;
}
#-----------------------------

# add edge $n1 -> n2
sub link_nodes {
   my ( $self, $n1, $n2, $word_id, $stt, $text) = @_;

   my %edge = ( n1 => $n1, n2 => $n2, word_id => $word_id, stt => $stt, text => $text);
   push @{$self->{_e}}, \%edge;

   my $edge_id = (scalar @{$self->{_e}}) - 1;
   push @{$self->{_out}->{$n1}}, $edge_id;
   push @{$self->{_in}->{$n2}}, $edge_id;
}
#-----------------------------

sub unlink_edge {
   my ( $self, $e ) = @_;

   my $n1 = $self->{_e}->[$e]->{n1};
   my $n2 = $self->{_e}->[$e]->{n2};
   # edge_id $e is deleted from incommin and outgoing edge arrays 
   main::del_int($self->{_out}->{$n1}, $e);
   main::del_int($self->{_in}->{$n2}, $e);
   # edge is not removed from the array, but marked as deleted
   $self->{_e}->[$e]->{stt} = -1; 
}
#-----------------------------

sub unlink_edges {
   my ( $self, $edges_ref, $from, $to) = @_;
   print "Unlinking edges [".join(',', @$edges_ref[$from..$to])."]\n" if ($main::debug == 1);
   for(my $i=$from; $i<=$to; $i++) {  
      $self->unlink_edge($$edges_ref[$i]);
      }
}
#-----------------------------

sub unlink_alternatives {
   my ( $self ) = @_;

   for(my $i=0; $i<scalar @{$self->{_e}}; $i++) {
      if ($self->{_e}->[$i]->{stt} == 0) {
         $self->unlink_edge($i); # stt becomes -1
         }
      }
   }
#-----------------------------

sub get_worst_stt {
   my ( $self, $edges_ref, $from, $to) = @_;

   for(my $i=$from; $i<=$to; $i++) {
      return 0 if ($self->{_e}->[$$edges_ref[$i]]->{stt} == 0);
      # We should never use the path over deleted edge
      printf "Error: path over deleted edge %d\n", $$edges_ref[$i] if ($self->{_e}->[$$edges_ref[$i]]->{stt} == -1 && $main::debug == 1);
      }
   return 1;
}
#-----------------------------

sub set_stt {
   my ( $self, $edges_ref, $from, $to, $new_stt) = @_;

   for(my $i=$from; $i<=$to; $i++) {
      # next if ($self->{_e}->[$$edges_ref[$i]]->{stt} == -1); # skip deleted edges
      # We should never use the path over deleted edge
      printf "Error: path over deleted edge %d\n", $$edges_ref[$i] if ($self->{_e}->[$$edges_ref[$i]]->{stt} == -1 && $main::debug == 1);
      $self->{_e}->[$$edges_ref[$i]]->{stt} = $new_stt if ($self->{_e}->[$$edges_ref[$i]]->{stt} != -1);
      }
}
#-----------------------------

sub multi_income {
   my ( $self, $edges_ref, $ignore_edge ) = @_;

   for(my $i=0; $i<scalar @$edges_ref; $i++) {
      next if ($$edges_ref[$i] == $ignore_edge); 
      return 1 if ($self->{_e}->[$$edges_ref[$i]]->{word_id} != -1);
      }
   return 0;
}
#-----------------------------

my @merges; # global array
my %uniq_seq;

sub batch_merge {
   my ( $self ) = @_;

#   print "Batch merge ---- \n" if ($main::debug == 1);
#   print join(' : ', @merges)."\n";
#return;
   # make keep list from all 1-length edge sequences
   my %keep_list = ();
   for(my $i=0; $i<scalar @merges; $i+=2) {
      my @edges = @{$merges[$i]};      
      $keep_list{$edges[0]} = 1 if ((scalar @edges) == 1);
      }

   #print "Keep list: ".join(' ', keys %keep_list)."\n" if ((scalar (keys %keep_list)) > 0) ;

   # Create new aggregated edges that replaces paths over the edge sequence 
   # Processed before unlinking to avoid iteration over broken paths
   for(my $i=0; $i<scalar @merges; $i+=2) {
      my @edges = @{$merges[$i]};
      my $num_text = $merges[$i+1];
      next if ((scalar @edges) == 1); # process edge sequences having lengths > 1

      my $from = 0;
      my $to = (scalar @edges)-1;
      my $n1 = $self->{_e}->[$edges[$from]]->{n1};
      my $n2 = $self->{_e}->[$edges[$to]]->{n2};
      # status = 1 if and only if all members of the sequence have status = 1
      #print "lastin=$$edges_ref[$lastin], lastout=$$edges_ref[$lastout]\n";
      my $stt = $self->get_worst_stt(\@edges, $from, $to);
      $self->set_stt(\@edges, $from, $to, 0) if ($stt == 1); # if this path is the best one, reset components to stt=0 
      $self->link_nodes($n1, $n2, -1, $stt, $num_text);
      print "Linking nodes ".$n1."-".$n2." stt=".$stt." text=".$num_text."\n" if ($main::debug == 1);
      }

   for(my $i=0; $i<scalar @merges; $i+=2) {
      my @edges = @{$merges[$i]};
      my $num_text = $merges[$i+1];
      next if ((scalar @edges) == 1); # process edge sequences having lengths > 1

      # Scan edges forward to find the first node having another incomming edge $firstin
      my $from = 0;
      my $to = (scalar @edges)-1;
      my $firstin = $to;
      for(my $j=$from; $j<$to; $j++) {
         my $n2 = $self->{_e}->[$edges[$j]]->{n2}; # $edges[$j] may be already deleted on another path
         if ($self->multi_income($self->{_in}->{$n2}, $edges[$j]) == 1) { # reikia neskaiciuoti naujø ryðiø
#         if ((scalar @{$self->{_in}->{$n2}}) > 1) { # reikia neskaiciuoti naujø ryðiø
            $firstin = $j;
            last;
            }
         }
      print "firstin = $edges[$firstin]\n" if ($main::debug == 1);
      # unlink [$from, $firstin] edges unless it is present in the $keep_list 
      #$self->unlink_edges(\@edges, $from, $firstin); # probably enough

      my $unlink_txt = '';
      for(my $j=$from; $j<=$firstin; $j++) {
         if (!defined $keep_list{$edges[$j]}) {  
            $self->unlink_edge($edges[$j]);
            $unlink_txt .= $self->{_e}->[$edges[$j]]->{n1}."-".$self->{_e}->[$edges[$j]]->{n2}." ";
            }
         else {
            print "Keep list: ".join(' ', keys %keep_list)."\n" if ((scalar (keys %keep_list)) > 0 && $main::debug == 1);
            }
         }
      print "Unlinking: ".$unlink_txt."\n" if ($main::debug == 1);
      }
   }
#-----------------------------

sub show_search_state {
   my ( $self, $edges_ref, $words_ref, $flag ) = @_;

   print $flag." ";
   print "Edges: [";
   for(my $i=0; $i<scalar @$edges_ref; $i++) {
      print " " if ($i > 0);
      print $$edges_ref[$i].($self->{_e}->[$$edges_ref[$i]]->{stt} == 0 ? '-' : '+');
      }
   print "] ";
   print "Nodes: [";
   print $self->{_e}->[$$edges_ref[0]]->{n1};
   for(my $i=0; $i<scalar @$edges_ref; $i++) {
      print " ".$self->{_e}->[$$edges_ref[$i]]->{n2};
      }
   print "] ";
   print "Text: '";
   for(my $i=0; $i<scalar @$words_ref; $i++) {
      print " " if ($i > 0);
      print $$words_ref[$i];
      printf "(%.2f)", $self->edge2time($$edges_ref[$i]) if ($$words_ref[$i] eq '<eps>'); 
      }
   print "'\n";
   }
#-----------------------------

sub extend_chain {
   my ( $self, $edges_ref, $words_ref, $codes_ref ) = @_;

   #print "Edges top: [".join(' ', @$edges_ref)."] Text: ".join(' ', @$words_ref)."\n" if ($main::debug == 1);
   $self->show_search_state($edges_ref, $words_ref, 'T') if ($main::debug == 1);
   
   #my $proc_done = 0; # becomes 1 to inhibit multiple recursion calls (2 tûkst. 2 m. + 2 tûkst. 2 metø)
   my $rc;
   my $n_curr = $self->{_e}->[$$edges_ref[-1]]->{n2};            # $n_curr is second node of the last edge in the array
   my $may_teminate_here = 0;
   $may_teminate_here = 1 if (!defined $self->{_out}->{$n_curr} || (scalar @{$self->{_out}->{$n_curr}}) == 0);
   # print "n_curr = $n_curr, out_edges = [".join(' ',@{$self->{_out}->{$n_curr}})."]\n";
   # iterate over all edge_ids in the out array of $n_curr
   # Iterating in reverse order so that edge deletions and additions in _out list won't affect search
   # this is no longer important as deletions are performed after recursive search returns
   if (defined $self->{_out}->{$n_curr}) {
      for(my $i=(scalar @{$self->{_out}->{$n_curr}})-1; $i>=0; $i--) {   
         my $edge_id = $self->{_out}->{$n_curr}->[$i];
         next if ($self->{_e}->[$edge_id]->{stt} == -1); # skip deleted edge, we should normally never get this here
         next if ($self->{_e}->[$edge_id]->{word_id} == -1); # skip added edge
         my $word = main::int2word($self->{_e}->[$edge_id]->{word_id});
         my $code = main::get_word_code(\$word);    
         $code = '-' if ($code eq 's' && $self->edge2time($edge_id) > $main::max_intra_sil); # <eps>
         push @$edges_ref, $edge_id;
         push @$words_ref, $word;
         push @$codes_ref, $code;
         if ($code eq '-') {        # chain ends, recursion terminates
            $may_teminate_here = 1;
            $self->show_search_state($edges_ref, $words_ref, 'N') if ($main::debug == 1);
            }
         else {
           $rc = $self->extend_chain($edges_ref, $words_ref, $codes_ref); # chain is extended, recursion continues
           }
         pop @$edges_ref;
         pop @$words_ref;
         pop @$codes_ref;
         }
      }

   if ($may_teminate_here == 1 && (!defined $rc || $rc ne 'M') ) {
#   if ($may_teminate_here == 1 && (scalar @$words_ref) > 1 && (!defined $rc || $rc ne 'M') ) {
      # process chain
      # return array of triplets <from, to, number>
      # triplets may be overlapping 
      #$proc_done = 1; # never get here again in the same $i loop with a different stop (-) word (e.g. m. or metø)
      $self->show_search_state($edges_ref, $words_ref, 'A') if ($main::debug == 1);
      my @result = ();
      # analyze_chunk may itself return a $multi_case condition
      my $multi_case = main::analyze_chunk($words_ref, $codes_ref, \@result, 0); 
      #printf "Result size=%d\n", scalar @result;
      if ($multi_case == 0) { # single explanation of a chunk
         for(my $j=0; $j<scalar @result; $j+=3) {
            my $from      = $result[$j];
            my $to        = $result[$j+1];
            my $num_text  = $result[$j+2];
            # extend $to to the end of a sequence if <eps> is the last remaining element
            $to++ if ($to+1 == (scalar @$edges_ref) - 1 && $$codes_ref[$to+1] eq 's');

            # v2. memorize edge sequence to be merged later & text
            # may contain duplicates
            my $key = join(' ' , @$edges_ref[$from..$to]);
            if (defined $uniq_seq{$key}) {
               if ($uniq_seq{$key} ne $num_text) {
                  print STDERR "$0: Chain ".$key." has different numeric representations '".$uniq_seq{$key}."' and '".$num_text."'\n";
                  }
               }
            else {
               $uniq_seq{$key} = $num_text;
               push @merges, ([@$edges_ref[$from..$to]], $num_text);    
               print "Merge [$key] -> $num_text\n" if ($main::debug == 1);
               }            
            }
         return 'S'; # single explanantion
         }
      else { # multiple explanations for a chunk, we are not joining anything
         print "Multi case\n" if ($main::debug == 1);
         return 'M'; # multiple explanations
         }
      }
   return $rc;
}
#-----------------------------

sub sort_edges_by_time {
   my ( $self ) = @_;
   my @keys = sort {
          $self->{_n_tm}->{$self->{_e}->[$a]->{n1}} <=> $self->{_n_tm}->{$self->{_e}->[$b]->{n1}}
                               ||
          $self->{_n_tm}->{$self->{_e}->[$a]->{n2}} <=> $self->{_n_tm}->{$self->{_e}->[$b]->{n2}} 
                               ||
          $self->{_e}->[$a]->{word_id} <=> $self->{_e}->[$b]->{word_id}
        } keys @{$self->{_e}};
   return \@keys;
   }
#-----------------------------

sub purge_del_edges {
   my ( $self, $keys_ref ) = @_;

   my $j = 0;
   for(my $i=0; $i<scalar @$keys_ref; $i++) {
      if ($self->{_e}->[$$keys_ref[$i]]->{stt} != -1) {
         $$keys_ref[$j] = $$keys_ref[$i] if ($i > $j); 
         $j++;
         }
      }
   splice (@$keys_ref, $j);
   return $keys_ref;
   }
#-----------------------------

# procedûros tikslas atstatyti skaièius pagal skaitinius komponentus
sub connect_nums {
   my ( $self ) = @_;

   my @keys = @{$self->sort_edges_by_time()};   
   for(my $i=0; $i<scalar @keys; $i++) { # $keys[$i] is edge index
      next if ($self->{_e}->[$keys[$i]]->{stt} == -1); # skip deleted edge
      next if ($self->{_e}->[$keys[$i]]->{word_id} == -1); # skip added edge
      my $word = main::int2word($self->{_e}->[$keys[$i]]->{word_id});
      my $code = main::get_word_code(\$word);    
      next if ($code eq '-');  
      next if ($code eq 's' && $self->edge2time($keys[$i]) > $main::max_intra_sil); # <eps>

      # accept as starting point any edge that either
      # has no parents 
      # has at least one '-' input edge 
      # has at least one 's' input edge longer than frames (<eps>)
      my $n1 = $self->{_e}->[$keys[$i]]->{n1};
      my $j;
      my $nb_parents = 0;
      if (defined $self->{_in}->{$n1}) {
         $nb_parents = scalar @{$self->{_in}->{$n1}};
         for($j=0; $j<$nb_parents; $j++) {
            my $edge_id = $self->{_in}->{$n1}->[$j];
            next if ($self->{_e}->[$edge_id]->{stt} == -1); # skip deleted edge
            next if ($self->{_e}->[$edge_id]->{word_id} == -1); # skip added edge
            my $parent_word = main::int2word($self->{_e}->[$edge_id]->{word_id});
            my $parent_code = main::get_word_code(\$parent_word);    
            last if ($parent_code eq '-');
            last if ($parent_code eq 's' && $self->edge2time($edge_id) > $main::max_intra_sil); # <eps>
            }
         }
      if (!defined $self->{_in}->{$n1} || $j < $nb_parents) {
         my @ch_edges = ($keys[$i]);
         my @ch_words = ($word);
         my @ch_codes = ($code);
         @merges = (); # empty @merges array
         %uniq_seq = ();
         $self->extend_chain(\@ch_edges, \@ch_words, \@ch_codes);
         $self->batch_merge();
         }
      }
}
#-----------------------------

# procedûros tikslas sumaþinti briaunø/þodþiø skaièiø redaktoriui
sub reduce_out {
   my ( $self ) = @_;

   # surûðiuojame briaunas pagal laikà ir þodá 
   # sort edges numerically
#  print join(' ', keys @{$self->{_e}} )."\n\n";
   my $key_ref = $self->sort_edges_by_time();   
   return $key_ref if ($main::debug == 1);
   my @keys = @{$self->purge_del_edges($key_ref)};

   # naikiname identiðkas briaunas
   my $j = 0;
   for(my $i=1; $i<scalar @keys; $i++) {
      if (!(
      $self->{_n_tm}->{$self->{_e}->[$keys[$i]]->{n1}} == $self->{_n_tm}->{$self->{_e}->[$keys[$j]]->{n1}} &&
      $self->{_n_tm}->{$self->{_e}->[$keys[$i]]->{n2}} == $self->{_n_tm}->{$self->{_e}->[$keys[$j]]->{n2}} &&
      $self->{_e}->[$keys[$i]]->{word_id} == $self->{_e}->[$keys[$j]]->{word_id}
      )) {
         $j++;
         $keys[$j] = $keys[$i];
         }
      elsif ($self->{_e}->[$keys[$i]]->{stt} == 1) { # main path
         $keys[$j] = $keys[$i];
         }
      }
   splice (@keys, $j+1);

   # naikiname smarkiai persidengianèias briaunas
   for(my $i=0; $i<scalar @keys; $i++) {
      next if ($self->{_e}->[$keys[$i]]->{stt} == -1); # skip deleted edge
      my $li = $self->{_n_tm}->{$self->{_e}->[$keys[$i]]->{n2}} - $self->{_n_tm}->{$self->{_e}->[$keys[$i]]->{n1}}; # length of i-th edge
      my $ti;
      if ($self->{_e}->[$keys[$i]]->{word_id} == -1) { # new joint number
         $ti = $self->{_e}->[$keys[$i]]->{text};
         }
      else { 
         $ti = main::int2word($self->{_e}->[$keys[$i]]->{word_id});
         $ti = $self->{_e}->[$keys[$i]]->{p2g} if (main::is_unk($ti));   
         }  
      my $ti_mod = $ti =~ s/[$qou$qcu]//gr;    
      $ti_mod = $1 if ($ti_mod =~ m/^(\d+)(-[[:alpha:]]+)?$/); # strip -osios if present
      $ti_mod = lc($ti_mod);
      $ti_mod = 'd.' if ($ti_mod =~ m/^dien(?:a|os|ai|$::a_nosine|os|$::u_nosine|oms|as|omis)$/);
      $ti_mod = 'm.' if ($ti_mod =~ m/^met(?:ai|$::u_nosine|ams|us|ais|uose)$/);
      my $tj;
      my $tj_mod;
      for(my $j=$i+1; $j<scalar @keys; $j++) {
         next if ($self->{_e}->[$keys[$j]]->{stt} == -1); # skip deleted edge
         last if ($self->{_n_tm}->{$self->{_e}->[$keys[$j]]->{n1}} >= $self->{_n_tm}->{$self->{_e}->[$keys[$i]]->{n2}});
         my $ident = 0;
         # Cannot compare word_ids, because of <unk>'s and joined edges (word_id == -1)
         #if ($self->{_e}->[$keys[$i]]->{word_id} == $self->{_e}->[$keys[$j]]->{word_id}) {  # identical word ids
         #   $ident = 1;
         #   }
         #else {
            if ($self->{_e}->[$keys[$j]]->{word_id} == -1) { # new joint number
               $tj = $self->{_e}->[$keys[$j]]->{text};
               }
            else {  
               $tj = main::int2word($self->{_e}->[$keys[$j]]->{word_id});
               $tj = $self->{_e}->[$keys[$j]]->{p2g} if (main::is_unk($tj));  
               }   
            $tj_mod = $tj =~ s/[$qou$qcu]//gr;
            $tj_mod = $1 if ($tj_mod =~ m/^(\d+)(-[[:alpha:]]+)?$/); # strip -osios if present
            $tj_mod = lc($tj_mod);
            $tj_mod = 'd.' if ($tj_mod =~ m/^dien(?:a|os|ai|$::a_nosine|os|$::u_nosine|oms|as|omis)$/);
            $tj_mod = 'm.' if ($tj_mod =~ m/^met(?:ai|$::u_nosine|ams|us|ais|uose)$/);

            #print "'$ti' '$tj' -> '$ti_mod' '$tj_mod'\n";
            $ident = 1 if ($ti_mod eq $tj_mod);
         #   }
         if ($ident == 1) {
            #print "'$ti' '$tj' -> '$ti_mod' '$tj_mod'\n";
       
            my $d1 = abs($self->{_n_tm}->{$self->{_e}->[$keys[$j]]->{n1}} - $self->{_n_tm}->{$self->{_e}->[$keys[$i]]->{n1}});
            my $d2 = abs($self->{_n_tm}->{$self->{_e}->[$keys[$j]]->{n2}} - $self->{_n_tm}->{$self->{_e}->[$keys[$i]]->{n2}});
            my $lj = $self->{_n_tm}->{$self->{_e}->[$keys[$j]]->{n2}} - $self->{_n_tm}->{$self->{_e}->[$keys[$j]]->{n1}};

            if ($d1 <= 2 || $d2 <= 2 || (List::Util::min($li, $lj) / List::Util::max($d1, $d2) > 4.0)) {
               # if you need to choose among two alternative hypotheses (stt=0) keep the lowercased one 
               if ($self->{_e}->[$keys[$j]]->{stt} == 0 && $self->{_e}->[$keys[$i]]->{stt} == 1 || 
                   $self->{_e}->[$keys[$j]]->{stt} == 0 && $self->{_e}->[$keys[$i]]->{stt} == 0 && (!defined($tj) || $tj !~ /^[[:lower:]]+$/) 
                  ) {
                  splice(@keys, $j, 1);
                  $j--;
                  }
               else {
                  splice(@keys, $i, 1);
                  $i--;
                  last;
                  }
               }
            }
         }

      }

   return \@keys;
#   $self->print_sorted(\@keys, 'STOI');
   }
#-----------------------------

sub label_best_v3 {
   my ( $self, $best ) = @_;

   my $n = $self->find_entry_node();
   my $n_b = $best->find_entry_node();
   my $ex_b = $best->find_exit_node();

   my $rc = label_best_loop_recursive($self ,$best, $n, $n_b, $ex_b, 0);

   if ($rc == 0) { # no path of L2 was found in L1
      print STDERR "Error: Mismatching full (L1) and best-path (L2) lattices\n";
      exit -1;
      }
   }
#-----------------------------

sub label_best_loop_recursive {
   my ( $self, $best, $n, $n_b, $ex_b, $depth ) = @_;

   my @sel_edges = ();

   # Travel through the best path
   while ($n_b != $ex_b) {
      my $e_b = $best->{_out}->{$n_b}->[0];
      my $edge_b = $best->{_e}->[$e_b];

      # iterate over all outgoing edges of the node $n in the complete graph
      # to find an edge matching $edge_b
      my $edges_ref = $self->{_out}->{$n};
      my $match_i = -1;
      my $match_cnt = 0;
      for(my $i=0; $i<scalar @$edges_ref; $i++) {
         my $edge = $self->{_e}->[$$edges_ref[$i]];
         if ($edge->{word_id} == $edge_b->{word_id} && $edge->{digest} eq $edge_b->{digest}) {
            $match_i = $i;
            $match_cnt++;
            }
         } 

      if ($match_cnt == 1) { # ordinary case, single path available, avoid recursion but loop
         push @sel_edges, $$edges_ref[$match_i];
         #$self->{_e}->[$$edges_ref[$match_i]]->{stt} = 1; # set status
         $n = $self->{_e}->[$$edges_ref[$match_i]]->{n2};
         $n_b = $edge_b->{n2};
         }
      # no match found
      elsif ($match_cnt == 0) { # no match found
         # if we are in recursion, we need to backtrack
         return 0 if ($depth > 0);
         # if we are not, we need to terminate on error
         my $fail_edge = '['.$edge_b->{n1}.' '.$edge_b->{n2}.'] ['.$best->node2time($edge_b->{n1}).' '.$best->node2time($edge_b->{n2}).'] '.$edge_b->{word_id}.' '.$edge_b->{digest};
         print STDERR "Error: Mismatch comparing lattices '$self->{_name}'. Entry '".$fail_edge."' was not found in L1\n";
         exit -1;
         }
      # multiple matches found
      elsif ($match_cnt > 1) {
         # Print anomaly to LOG file
         my $fail_edge = '['.$edge_b->{n1}.' '.$edge_b->{n2}.'] ['.$best->node2time($edge_b->{n1}).' '.$best->node2time($edge_b->{n2}).'] '.$edge_b->{word_id}.' '.$edge_b->{digest};
         print STDERR "Warning: Mismatch comparing lattices '$self->{_name}'. Entry '".$fail_edge."' has multiple matches in L1\n";
         for(my $i=0; $i<scalar @$edges_ref; $i++) {
            my $edge = $self->{_e}->[$$edges_ref[$i]];
            if ($edge->{word_id} == $edge_b->{word_id}  && $edge->{digest} eq $edge_b->{digest}) {
               my $fail_edge = '['.$edge->{n1}.' '.$edge->{n2}.'] ['.$self->node2time($edge->{n1}).' '.$self->node2time($edge->{n2}).'] '.$edge->{word_id}.' '.$edge->{digest};
               print STDERR $fail_edge."\n";
               }
            }
         # recursively iterate through every path
         $match_cnt = 0;
         for(my $i=0; $i<scalar @$edges_ref; $i++) {
            my $edge = $self->{_e}->[$$edges_ref[$i]];
            if ($edge->{word_id} == $edge_b->{word_id}  && $edge->{digest} eq $edge_b->{digest}) {

               push @sel_edges, $$edges_ref[$i];
               # $self->{_e}->[$$edges_ref[$i]]->{stt} = 1; # set status (on the best path)
               $match_cnt++;
               my $rc = label_best_loop_recursive($self, $best, $edge->{n2}, $edge_b->{n2}, $ex_b, $depth+1);
               if ($rc == 0) {
                  pop @sel_edges;
                  #$self->{_e}->[$$edges_ref[$i]]->{stt} = 0; # restore status (not on the best path)
                  $match_cnt--;
                  }
               }
            } 
         if ($match_cnt == 0) { # no match found, dead end for the recursion
            return 0 if ($depth > 0);
            print STDERR "Error: Disambiguation among alternative paths has failed.\n";
            exit -1;
            }
         elsif ($match_cnt == 1) { # disambiguation succeeded
            return 1 if ($depth > 0);
            last;
            }
         elsif ($match_cnt > 1) { # multiple matches found
            print STDERR "Error: Disambiguation among alternative paths has failed.\n";
            exit -1;
            }
         }
      }

   foreach ( @sel_edges ) {
      $self->{_e}->[$_]->{stt} = 1; # set status
      }

   return 1;
   }
#-----------------------------

sub label_phones {
   my ( $self, $phones ) = @_;

   # assuming the same number of edges
   for (my $i=0; $i<scalar @{$self->{_e}}; $i++) {
      # space-separated phone list (non-BESI)
      $self->{_e}->[$i]->{ph} = $phones->{_e}->[$i]->{ph}; 
      # convert phone sequence to word (graphemes)
      # word-by-word conversion is not efficient
      # $self->{_e}->[$i]->{p2g} = phones2letters($phones->{_e}->[$i]->{ph}) if (is_unk(main::int2word($self->{_e}->[$i]->{word_id})))
      }

   # collect all unknown words into a list so that they can be sent in a single batch to the phone2word service
   my %unk_hash = ();
   my @unk_list = ();
   for (my $i=0; $i<scalar @{$self->{_e}}; $i++) {
      next if ($self->{_e}->[$i]->{stt} == -1); # skip words marked for deletion
      if (main::is_unk(main::int2word($self->{_e}->[$i]->{word_id}))) {
         if (!defined($unk_hash{$self->{_e}->[$i]->{ph}})) {
            push @unk_list, {'phones' => $self->{_e}->[$i]->{ph}};
            $unk_hash{$self->{_e}->[$i]->{ph}} = 1;
            }
         }
      }
   return if (scalar @unk_list == 0);

   main::phones2letters_list(\@unk_list);
   foreach (@unk_list) {
      $unk_hash{$_->{phones}} = $_->{word};
      # print "$_->{phones} -> $_->{word}\n"; # debug
      }

   for (my $i=0; $i<scalar @{$self->{_e}}; $i++) {
      next if ($self->{_e}->[$i]->{stt} == -1);
      if (main::is_unk(main::int2word($self->{_e}->[$i]->{word_id}))) {
         $self->{_e}->[$i]->{p2g} = $unk_hash{$self->{_e}->[$i]->{ph}};
         }
      }
   }
#-----------------------------

sub print_sorted {
   # keys_ref is list of $self->{_e} indexes in some order
   my ( $self, $keys_ref, $options ) = @_;

   if ($options =~ m/^[FSNTIWPODE]+$/ ) {
      if ($options =~ m/F/) {
         print "$self->{_name}";
         print " " if ($options =~ m/E/); # 2020.11.18 jei norime, kad butu, kaip originale
         print "\n";
         }

      my @P = split('', $options);
      for (my $i=0; $i<scalar @{$keys_ref}; $i++) {
         my $edge = $self->{_e}->[$$keys_ref[$i]]; 
         # next if ($edge->{stt} == -1); # skip deleted edge
         my $line = '';
         for(my $j=0; $j<scalar @P; $j++) {
            if ($j > 0 && $P[$j] !~ m/^[FE]$/) {
               if ($options =~ m/E/) {
                  $line .= "\t" # originalus lattice skirtukas
                  }
               else {
                  $line .= ' '; 
                  }
               }
            $line .= $edge->{stt} if ($P[$j] eq 'S'); # Status
            $line .= $edge->{n1}.' '.$edge->{n2} if ($P[$j] eq 'N'); # Nodes
            $line .= $self->node2time($edge->{n1}).' '.$self->node2time($edge->{n2}) if ($P[$j] eq 'T'); # Timing
            $line .= $edge->{word_id} if ($P[$j] eq 'I'); # Word ID
            if ($P[$j] eq 'W') { # Word text
               if ($edge->{word_id} == -1) { # new joint number
                  $line .= $edge->{text};
                  }
               else {
                  $line .= main::int2word($edge->{word_id});
                  }
               }
            $line .= '<'.$edge->{ph}.'>' if ($P[$j] eq 'P' && defined $edge->{ph}); # Phone sequence
            if ($P[$j] eq 'O') { # output for editor
               if ($edge->{word_id} == -1) { # new joint number
                  my $word = $edge->{text};
                  # pasaliname priedus 1-as 2-u
                  $word = $main::rev_map_digits_N{$word} if (defined $main::rev_map_digits_N{$word});
                  $line .= $word;
                  }
               elsif (!main::is_unk(main::int2word($edge->{word_id}))) {
                  my $word = main::int2word($edge->{word_id});
                  # jei pasaliname jungiamuosius bruksnius tarp zodziu - redaktoriuje jie sulimpa
                  # pasaliname priedus 1-as 2-u
                  $word = $main::rev_map_digits_N{$word} if (defined $main::rev_map_digits_N{$word});
                  $line .= $word;
                  }
               else {
                  my $word = $edge->{p2g}; # _N_ i kampinius skliaustus nededam
                  $word = '<'.$word.'>' if (main::int2word($edge->{word_id}) eq '<unk>');
                  $line .= $word;
                  }
               }
            if ($P[$j] eq 'D') { # digest / used for phone bondary shift
               $line .= $edge->{digest};
               }
            }
         printf("%s\n", $line);
         }
      if ($options =~ m/E/) { # end lattice so that it is exactly the same as original
         my $edge = $self->{_e}->[$$keys_ref[-1]];
         printf("%d\n\n", $edge->{n2});
         }
      }
   else {
      print "$self->{_name}\n"; 
      for my $node ( sort {$a <=> $b} keys %{$self->{_n_tm}} ) { 
         print "$node $self->{_n_tm}->{$node}\n";
         }
      print "\n"; 

      for my $node ( sort {$a <=> $b} keys %{$self->{_out}} ) { 
         print "Node = $node:\n";
         if (defined $self->{_in}->{$node}) {
            print "InArcs ".join(' ', @{$self->{_in}->{$node}})."\n"; 
            }
         print "OutArcs ".join(' ', @{$self->{_out}->{$node}})."\n";
         }
      print "\n"; 

      for (my $i=0; $i<scalar @{$self->{_e}}; $i++) {
         print "$i: $self->{_e}->[$i]\n";
         }
      print "\n"; 
   }
}
#-----------------------------

# print in default order
sub print {
   my ( $self, $options ) = @_;

   $self->print_sorted([keys @{$self->{_e}}], $options);
}
#-----------------------------

# tests if there are >1 "best" hypotheses
sub test_single_best {
   my ( $self, $id ) = @_;

   my @keys = @{$self->sort_edges_by_time()};   
   for(my $i=0; $i<scalar @keys; $i++) {
      next if ($self->{_e}->[$keys[$i]]->{stt} != 1); # look only for best edges
      my $i_end_frame = $self->{_n_tm}->{$self->{_e}->[$keys[$i]]->{n2}};
      for(my $j=$i+1; $j<scalar @keys; $j++) {
         next if ($self->{_e}->[$keys[$j]]->{stt} != 1); # look only for best edges
         my $j_beg_frame = $self->{_n_tm}->{$self->{_e}->[$keys[$j]]->{n1}};
         last if ($j_beg_frame == $i_end_frame);
         # overlap
         if ($j_beg_frame < $i_end_frame) {
            printf STDERR "$0: Lattice no. %d has best path overlap at [%s, %s]\n", $id, $self->frame2time($j_beg_frame), $self->frame2time($i_end_frame);
            }
         # skip
         else {
            printf STDERR "$0: Lattice no. %d has no best path over [%s, %s]\n", $id, $self->frame2time($i_end_frame), $self->frame2time($j_beg_frame);
            last;
            }
         }
      }
   
   # Count broken segments (nodes without either incoming or outgoing edges)
   # Nodes without both incoming and outgoing edges are ok
   my $cnt = 0;
   foreach my $n (keys %{$self->{_n_tm}}) {
      next if (!defined $self->{_in}->{$n}); # skip entry node
      next if (!defined $self->{_out}->{$n}); # skip exit node
      my $in_edges_cnt = scalar @{$self->{_in}->{$n}};
      my $out_edges_cnt = scalar @{$self->{_out}->{$n}};
      $cnt++ if (($in_edges_cnt > 0 && $out_edges_cnt == 0) || 
                 ($in_edges_cnt == 0 && $out_edges_cnt > 0));
      }
   printf STDERR "$0: Lattice no. %d has %d broken segments\n", $id, $cnt if ($cnt > 0);
   }
#-----------------------------


1;