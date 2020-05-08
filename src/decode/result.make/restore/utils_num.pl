#!/usr/bin/perl

#XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
#X                                                                  X
#X utils_num.pl, Copyright(C) Gailius Raðkinis, 2020                X
#X                                                                  X
#XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX


# Script that restores numbers from recognizer's output 
# 7 tûkst. 6 100 20 1 -> 7621 

use strict;
use warnings;
use File::Basename;

binmode(STDIN, ":utf8");
binmode(STDOUT, ":utf8");

my $dirname = dirname(__FILE__);
require "$dirname/lt.pl";
 
my @order = ("t\Q$::u_ilgoji\Ekst\.", "mln\.", "mlrd\.", "trln\.");

###############################################################################
#
# NUM TO TEXT CONVERSION
#
###############################################################################

# Converts number into text
# reiketu papildyti vns. dgs. vienareiksminimu 100-as,o 2 100-ai,ø
sub int2text {
   my ( $num ) = @_;
   my $o;
   my $i;

   #print "$num\n";
   my @digits = split('', reverse $num); 
   my $text = '';
   my $len = scalar @digits;

   # +/-
   if ($digits[-1] =~ m/([-+])/ ) {
      $text .= $1.' ';
      $len--;
      }

   for ($i=0; $i<$len; $i+=3 ) {
      my $tripl = '';
      if ($i+1 < $len) { # [..00 - ..99]
         if ($digits[$i+1] == 0) { # [..00 - ..09]
            $tripl = ($digits[$i] == 0 ? '' : $digits[$i]);
            }
         elsif ($digits[$i+1] == 1) { # [..10 - ..19]
            $tripl = $digits[$i+1].$digits[$i];
            }
         else { # [..20 - ..99]
            $tripl = $digits[$i+1].'0'.($digits[$i] == 0 ? '' : ' '.$digits[$i]); # 50 9
            }
         }
      else { # [0-9]
         $tripl = $digits[$i];
         }
      if ($i+2 < $len) { # simtai
         my $simt  = ($digits[$i+2] == 0 ? '' : ($digits[$i+2] == 1 ? '' : $digits[$i+2].' ').'100 ');
         $tripl = (($digits[$i] == 0 && $digits[$i+1] == 0) ? $simt : $simt.$tripl);
         }
      {
      use integer;
      $o = (($i / 3) - 1) % scalar @order; # order of magnitude
      }
      if ($o < 0) {
         $text = $tripl;
         }
      elsif ($tripl eq '') {
         }
      elsif ($tripl eq '1') {
         $text = $order[$o].' '.$text;
         }
      else {
         $text = $tripl.' '.$order[$o].' '.$text;
         }
      }

   return($text);
   }

###############################################################################
#
# TEXT TO NUM CONVERSION
#
###############################################################################

# get a character code for a word
sub get_word_code {
   my ( $word_ref ) = @_;

   if ($$word_ref =~ m/^<eps>$/) { # <eps>
      return 's'; 
      }
   if ($$word_ref =~ m/^100(?:-[[:alpha:]]+)?$/) { # 100, 100-aisiais
      return 'h'; 
      }
   if ($$word_ref =~ m/^[2-9]0(?:-[[:alpha:]]+)?$/) { # 20, ..., 90
#   if ($$word_ref =~ m/^(?:[2-9]0(?:-[[:alpha:]]+)?|(?:dvi|tris|keturias|penkias|ðeðias|septynias|aðtuonias|devynias)deðimt)$/ ) { # 20, ..., 90
      return 't'; 
      }
   if ($$word_ref =~ m/^1[0-9](?:-[[:alpha:]]+)?$/) { # 10-19
      return 'e'; 
      }
   if ($$word_ref =~ m/^[2-9](?:-[[:alpha:]]+)?$/) { # 2-9
      return 'u'; 
      }
   if ($$word_ref =~ m/^1(?:-[[:alpha:]]+)?$/) { # 1
      return '0'; 
      }
   for (my $j=0; $j<scalar @order; $j++) {
      return ($j+1) if ($$word_ref =~ m/^\Q$order[$j]\E$/);
      }
   return '-';
   }


# Returns 1 if chunk is broken at $pos]
# $next_pos is usually $pos+1, or $pos+2 is <eps> is inserted
sub is_num_break {
   my ( $types_ref, $from, $pos ) = @_;

   # $types_ref - sequence of num codes
   # $from      - left boundary of the current sub-sequence in $types_ref
   # $pos       - break position to test

   # assuming there aren't consecutive <eps>'es
   my $next_pos = $pos+1;
   $next_pos++ if ($next_pos < (scalar @$types_ref) && $$types_ref[$next_pos] eq 's'); # <eps>
   return 1 if ($next_pos >= scalar @$types_ref); # sequence end has been reached

   if ($$types_ref[$pos] eq '0') { # 1
      return 1 if (!($pos > $from && $$types_ref[$next_pos] =~ m/^[1234]$/));
      }
   elsif ($$types_ref[$pos] eq 'u') { # 2-9
      my $prev_pos = $pos-1;
      $prev_pos-- if ($prev_pos > $from && $$types_ref[$prev_pos] eq 's'); # <eps>
      return 1 if (!($$types_ref[$next_pos] eq 'h' && ($pos == $from || $prev_pos >= $from && $$types_ref[$prev_pos] =~ m/^[1234]$/) || 
                     $$types_ref[$next_pos] =~ m/^[1234]$/));
      }
   elsif ($$types_ref[$pos] eq 'e') { # 10-19
      return 1 if (!($$types_ref[$next_pos] =~ m/^[1234]$/));
      }
   elsif ($$types_ref[$pos] eq 't') { # 20-90
      return 1 if (!($$types_ref[$next_pos] =~ m/^[12340u]$/));
      }
   elsif ($$types_ref[$pos] eq 'h') { # 100
      return 1 if (!($$types_ref[$next_pos] =~ m/^[12340uet]$/));  # visi
      }
   elsif ($$types_ref[$pos] =~ m/^[1-4]$/) { # tûkst, mln, mlrd, trln
      return 1 if (!(($$types_ref[$next_pos] =~ m/^[0ueth]$/ || 
                      $$types_ref[$next_pos] =~ m/^[123]$/ && $$types_ref[$pos] > $$types_ref[$next_pos]))); # mln. tûkst.
      }
   return 0;
}

# Input is a close ended [] chunk consisting of non '-' tokens 
# Returns an array of triplets <from, to, number>
# about how the chunk can be covered by numbers
# Triplets may be complementary (one number is following another)
# or competitive 


sub analyze_chunk {
   my ( $words_ref, $types_ref, $res_ref, $debug ) = @_;
   # $res_ref is assumed empty

   my $nb_words = scalar @$words_ref;
   my $multi_case = 0;
   for (my $from=0; $from < $nb_words; $from++) {
      # return the longest possible numerical chain starting from $from
      next if ($$types_ref[$from] =~ m/^[-s]$/); # should not be present in a chunk
      #print "from=$from\n";

      my $smnum = 0; # small num 1-999
      my $smfrom = $from;
      my $osios;
      my $to = -1;
      my @numlist = ();
      for (my $pos = $from; $pos < $nb_words; $pos++) { # for every word
         #print "pos=$pos\n";
         next if ($$types_ref[$from] =~ m/^s$/); # skip <eps>
         $$words_ref[$pos] =~ m/^(\d+)(-[[:alpha:]]+)?$/ if ($$types_ref[$pos] !~ m/^[1234]$/); # strip -osios if present
         my $word = $1;
         $osios = undef;
         #$add_order = undef;  
         # number chan is interrupted by -osios
         if (defined $2) {
            $osios = $2;
            $to = $pos;
            }
         if ($$types_ref[$pos] =~ m/^[0uet]$/) { # 1, 2-9, 10-19, 20-90
            $smnum += $word; 
            }
         if ($$types_ref[$pos] eq 'h') { # 100
            $smnum = 1 if ($smnum == 0);
            $smnum = $smnum * 100; 
            }
         if ($$types_ref[$pos] =~ m/^[1234]$/) { # tûkst, mln, mlrd, trln
            $smnum = 1 if ($smnum == 0);
            # fiktyvus vienetas duoda  _from > _to
            my %num_ext = ( _num => $smnum, _from => $smfrom, _to => $pos-1, _ext => $osios, _ord => $$types_ref[$pos] );
            push @numlist, \%num_ext;

            $smnum = 0;
            $smfrom = $pos+1;
            }
         $to = $pos if (is_num_break($types_ref, $from, $pos) == 1);
         last if ($to != -1); 
         } # for pos

      # 0 
      if ($to == -1 && $nb_words == 1 && $$words_ref[0] eq '0') { 
         $to = 0;
         my %num_ext = ( _num => 0, _from => $smfrom, _to => $smfrom, _ord => 0 );
         push @numlist, \%num_ext;
         }
      # 1-999
      elsif ($smnum > 0) {
         my %num_ext = ( _num => $smnum, _from => $smfrom, _to => $to, _ext => $osios, _ord => 0 );
         push @numlist, \%num_ext;
         }

      # debug
      if ($debug == 1) {
         for(my $i=0; $i<scalar @numlist; $i++) {
            printf "%5d %3d %3d %8s %1d\n", $numlist[$i]->{_num}, $numlist[$i]->{_from}, $numlist[$i]->{_to}, 
                                            (defined $numlist[$i]->{_ext} ? $numlist[$i]->{_ext} : 'n/a'), $numlist[$i]->{_ord}
            }
         printf "\n";
         }
      # debug


      # Check if natural orders are reversed
      # 100 50 tûkst. mln. (150 tûkstanèiø milijonø)
      # 100 tûkst. 30 2 mln. (100 tûkstanèiø, 32 milijonai)
      my $rev_order = 0;
      for(my $i=1; $i<scalar @numlist; $i++) {
         if ($numlist[$i-1]->{_ord} < $numlist[$i]->{_ord}) {
            $rev_order = 1;
            }
         }
      $multi_case = 1 if ($rev_order == 1);

      # new number must extend beyond the previous
      next if ((scalar @$res_ref) > 0 && $numlist[-1]->{_to} <= $$res_ref[-2]);

      # 122 mln. 524 tûkst. -> do not convert numbers without units (1-999) into long integers 
      if ($rev_order == 1 || $numlist[-1]->{_ord} > 0) {
         # we must include these segments so that they are not deleted later
         for(my $i=0; $i<scalar @numlist; $i++) {
            push @$res_ref, ($numlist[$i]->{_from}, $numlist[$i]->{_to}, $numlist[$i]->{_num}) 
               if ($numlist[$i]->{_from} <= $numlist[$i]->{_to}); # ne fiktyvus 1
            push @$res_ref, ($numlist[$i]->{_to}+1, $numlist[$i]->{_to}+1, $$words_ref[$numlist[$i]->{_to}+1]) 
               if ($numlist[$i]->{_ord} > 0);                     # tûkst. mln.
            }
         }
      # 122 mln. 524 tûkst. 114 -> 122 524 114 
      else {
         my $num = 0;
         for(my $i=0; $i<scalar @numlist; $i++) {
            $num += $numlist[$i]->{_num} * (1000 ** $numlist[$i]->{_ord}); 
            }
         # long "pure" numbers should be made mored readable by inserting separators
         # e.g. 19850056456 -> 19_850_056_456
         if (length($num) >= 5) {
            $num = reverse $num;
            $num =~ s/...\K(?=.)/_/sg;
            $num = reverse $num;
            }
         # return 19-osios
         $num .= $osios if (defined $osios);
         push @$res_ref, ($numlist[0]->{_from}, $numlist[-1]->{_to}, $num); 
         }
       #print "from=$from to=$to num=$num\n";
      } # for from   

   # debug
   if ($debug == 1) {
      print '---'."\n";
      for(my $i=0; $i<scalar @$res_ref; $i+=3) {
         print 'from='.$$res_ref[$i].' to='.$$res_ref[$i+1].' num='.$$res_ref[$i+2]."\n";
         }
      # debug
      print '---'."\n";
      }
   #print "A: multicase=$multi_case\n";
   $multi_case = 1 if ($multi_case == 0 && multiple_num_hyp($res_ref) == 1);
   #print "B: multicase=$multi_case\n";
   return($multi_case);
   }

# Tests if there are overlapping hypotheses
# Input an array of triplets <from, to, number>
sub multiple_num_hyp {
   my ( $res_ref ) = @_;

   for(my $j=3; $j<scalar @$res_ref; $j+=3) {
      return 1 if ($$res_ref[$j-2] >= $$res_ref[$j]);
      }
   return 0;
   }


1
