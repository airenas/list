#!/usr/bin/perl 

use strict;
use warnings;
use Encode qw(decode encode);
use File::Basename;
use List::Util qw[min max];
#use Data::Dumper qw(Dumper);

my $qou = decode("utf-8", qq/\xe2\x80\x9e/); # „
my $qcu = decode("utf-8", qq/\xe2\x80\x9c/); # “

sub phones2letters {
   my ( $phstr ) = @_;

   $phstr =~ s/x/ch/g;
   $phstr =~ s/S/$::sh/g;
   $phstr =~ s/Z/$::zh/g;
   $phstr =~ s/E/$::eh/g;
   $phstr =~ s/N/n/g;
   $phstr =~ s/G/h/g;
   $phstr =~ s/ji(?=[ou])/j/g;
   $phstr =~ s/[\.:\'\"^]//g; 
   return($phstr);
   }


my $dirname = dirname(__FILE__);
require "$dirname/lat_map.pl";
require "$dirname//lt.pl";

package LatGraph;

sub new {
   my $class = shift;
   my $self = {
      _name => shift,
      _frameRate => shift,
      _startTime => shift,
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

# used to read L1 & L2 (frame-based)
sub add_edge {
   my ( $self, $n1, $n2, $word_id, $nb_frames ) = @_;

   my %edge = ( n1 => $n1, n2 => $n2, word_id => $word_id, stt => 0 );
   push @{$self->{_e}}, \%edge;
   my $edge_id = (scalar @{$self->{_e}}) - 1;

   $self->{_n_tm}->{$n2} = $self->{_n_tm}->{$n1} + $nb_frames; 

   $self->{_out}->{$n1} = () if (!defined $self->{_out}->{$n1});
   push @{$self->{_out}->{$n1}}, $edge_id;
   $self->{_in}->{$n2} = () if (!defined $self->{_in}->{$n2});
   push @{$self->{_in}->{$n2}}, $edge_id;
}

# used to read L3 (phone-based)
sub add_phones {
   my ( $self, $phone_ids ) = @_;

   my @C = split("_", $phone_ids);
   my $transc='';
   for (my $p = 0; $p <= $#C; $p++) {
      $transc = $transc.main::int2phone($C[$p], 1); # remove_pd_info
      }
#   $transc = '<'.$transc.'>';
   my %edge = ( ph => $transc );
   push @{$self->{_e}}, \%edge;
}

# procedûros tikslas sumaþinti briaunø/þodþiø skaièiø redaktoriui
sub reduce_out {
   my ( $self ) = @_;

   # surûðiuojame briaunas pagal laikà ir þodá (_out ir _in indeksai sugadinami) 
   # sort edges numerically
#  print join(' ', keys @{$self->{_e}} )."\n\n";
   my @keys = sort {
          $self->{_n_tm}->{$self->{_e}->[$a]->{n1}} <=> $self->{_n_tm}->{$self->{_e}->[$b]->{n1}}
                               ||
         $self->{_n_tm}->{$self->{_e}->[$a]->{n2}} <=> $self->{_n_tm}->{$self->{_e}->[$b]->{n2}} 
                               ||
         $self->{_e}->[$a]->{word_id} <=> $self->{_e}->[$b]->{word_id}
        } keys @{$self->{_e}};

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
      my $li = $self->{_n_tm}->{$self->{_e}->[$keys[$i]]->{n2}} - $self->{_n_tm}->{$self->{_e}->[$keys[$i]]->{n1}}; # length of i-th edge
      my $ti = main::int2word($self->{_e}->[$keys[$i]]->{word_id});
      $ti = main::phones2letters($self->{_e}->[$keys[$i]]->{ph}) if ($ti eq '<unk>');     
      my $ti_mod = $ti =~ s/[$qou$qcu]//gr;
      $ti_mod = lc($ti_mod);
      my $tj;
      my $tj_mod;
      for(my $j=$i+1; $j<scalar @keys; $j++) {
         last if ($self->{_n_tm}->{$self->{_e}->[$keys[$j]]->{n1}} >= $self->{_n_tm}->{$self->{_e}->[$keys[$i]]->{n2}});
         my $ident = 0;
         if ($self->{_e}->[$keys[$i]]->{word_id} == $self->{_e}->[$keys[$j]]->{word_id}) {  # identical word ids
            $ident = 1;
            }
         else {
            $tj = main::int2word($self->{_e}->[$keys[$j]]->{word_id});
            $tj = main::phones2letters($self->{_e}->[$keys[$j]]->{ph}) if ($tj eq '<unk>');     
            $tj_mod = $tj =~ s/[$qou$qcu]//gr;
            $tj_mod = lc($tj_mod);

            #print "'$ti' '$tj' -> '$ti_mod' '$tj_mod'\n";
            $ident = 1 if ($ti_mod eq $tj_mod);
            }
         if ($ident == 1) {
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

sub mk_timed_hash {
   my ( $self ) = @_;

   my %H;
   %{H} = ();
   for (my $i=0; $i<scalar @{$self->{_e}}; $i++) {
      my $edge = $self->{_e}->[$i];
#      $edge =~ m/^(\d+)\s(\d+)\s(\d+)\s(\d+)$/; 
#      if (!defined($1) || !defined($2) || !defined($3) || !defined($4)) {
#         print "'$edge'-> bad edge\n";
#         exit 0;
#         }
      # gali bûti keli vienodi $newkey. Ið jø iðsaugomas tik paskutinis
      my $newkey = $self->{_n_tm}->{$edge->{n1}}.' '.$self->{_n_tm}->{$edge->{n2}}.' '.$edge->{word_id};
      $H{$newkey} = $i;
      }

   return(\%H);
   }

sub label_best {
   my ( $self, $best ) = @_;

   # Construct hash 'time_from time_to word_id'
   my $H  = $self->mk_timed_hash();
   my $HB = $best->mk_timed_hash();

   # Gali paþymëti ne ta briaun1, kaip pagrindin3, nes H ir atsimenama tik viena ið keliø tapaèiø briaunø 
   for my $str ( keys %{ $HB } ) { 
      if (!defined $H->{$str}) { # look for an identical <from to int/text_id> triplet
         print STDERR "Error: Mismatch comparing lattices '$self->{_name}'. Entry '$str' was not found in L1\n";
         exit -1;
         }
      $self->{_e}->[$H->{$str}]->{stt} = 1; # replace terminal 0 by 1
      }
   }

sub label_phones {
   my ( $self, $phones ) = @_;

   # assuming the same number of edges
   for (my $i=0; $i<scalar @{$self->{_e}}; $i++) {
      $self->{_e}->[$i]->{ph} = $phones->{_e}->[$i]->{ph}; 
      }
   }

sub print_sorted {
   # keys_ref is list of $self->{_e} indexes in some order
   my ( $self, $keys_ref, $options ) = @_;

   if ($options =~ m/^[FSNTIWPO]+$/ ) {
      print "$self->{_name}\n" if ($options =~ m/F/); 

      my @P = split('', $options);
      for (my $i=0; $i<scalar @{$keys_ref}; $i++) {
         my $edge = $self->{_e}->[$$keys_ref[$i]]; 
         my $line = '';
         for(my $j=0; $j<scalar @P; $j++) {
            $line .= ' ' if ($j > 0);
            $line .= $edge->{stt} if ($P[$j] eq 'S'); # Status
            $line .= $edge->{n1}.' '.$edge->{n2} if ($P[$j] eq 'N'); # Nodes
            $line .= ($self->{_startTime} + $self->{_n_tm}->{$edge->{n1}}*$self->{_frameRate}).' '.($self->{_startTime} + $self->{_n_tm}->{$edge->{n2}}*$self->{_frameRate}) if ($P[$j] eq 'T'); # Timing
            $line .= $edge->{word_id} if ($P[$j] eq 'I'); # Word ID
            $line .= main::int2word($edge->{word_id}) if ($P[$j] eq 'W'); # Word text
            $line .= '<'.$edge->{ph}.'>' if ($P[$j] eq 'P' && defined $edge->{ph}); # Phone sequence
            if ($P[$j] eq 'O') { # output for editor
               if (main::int2word($edge->{word_id}) ne '<unk>') {
                  $line .= main::int2word($edge->{word_id});
                  }
               else {
                  $line .= '<'.main::phones2letters($edge->{ph}).'>';
                  }
               }
            }
         printf("%s\n", $line);
         }
      print "\n"; # newline between files
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

# print in default order
sub print {
   my ( $self, $options ) = @_;

   $self->print(keys @{$self->{_e}}, $options);
}


1;