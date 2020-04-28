#!/usr/bin/perl

# ./lat_restore.pl L1.lat L2.lat L3.lat words.txt phones.txt
# L1.lat Word-aligned lattice that has transition-id sequence
# L2.lat Best-path word-aligned lattice that has transition-id sequence
# L3.lat Word-aligned lattice that has phone-id sequence
# words.txt  -  word-id to word-string mapping
# phones.txt - phone-id to phone-string mapping

# ./lat2sym.pl /home/gailius/pcidisk/egs/g3_word/LAB/data/lang-BESI-1g/words.txt /home/gailius/pcidisk/egs/g3_word/LAB/data/lang-BESI-1g/phones.txt /home/gailius/pcidisk/egs/g2_noise/test_lat/tmplat.lat
# Lattice to text lattice conversion
# (after lattice-align-words and lattice-to-phone-lattice were performed)

use strict;
use warnings;
use File::Basename;
use lib dirname (__FILE__); # kad rastu .pm
use LatGraph;
#use Data::Dumper qw(Dumper);

my $dirname = dirname(__FILE__);
require "$dirname/lat_map.pl";
 
binmode(STDIN, ":utf8");
binmode(STDOUT, ":utf8");


###############################################################################
#                                                                             #
# MAIN PROGRAM                                                                #
#                                                                             #
###############################################################################

my $L1 = shift @ARGV; # multiple-path lattice with transition-ids
my $L2 = shift @ARGV; # single best path lattice with transition-ids
my $L3 = shift @ARGV; # multiple-path lattice with phone-ids
$main::symtab_w = shift @ARGV; # word symbol table
$main::symtab_p = shift @ARGV; # phone symbol table

if (!defined $L1 || !defined $L2 || !defined $L3 || !defined $main::symtab_w || !defined $main::symtab_p) {
   print STDERR "Usage: $0 L1.lat L2.lat L3.lat words.txt phones.txt\n" .
   exit(1);
   }

# Counts space separated entries in a string
sub count_spaces {
   my ( $line_ref ) = @_;
   return($$line_ref =~ tr/ //);
   }
#-----------------------------

sub count_tabs {
   my ( $line_ref ) = @_;
   return($$line_ref =~ tr/\t//);
   }
#-----------------------------

sub count_underscores {
   my ( $line_ref ) = @_;
   return($$line_ref =~ tr/_//);
   }
#-----------------------------

# $input_file may contain multiple lattices
sub read_lats {
   my ( $input_file, $lats_ref, $frameRate, $frameBased ) = @_;
   my @lines; # read entire file into array

   # Read input file
   my $rc = (open my $handle, '<:utf8', $input_file); 
   if (defined $rc && $rc != 0) {
      @lines = <$handle>;
      close $handle;
      }
   else {
      warn "$0: Couldn't open file $input_file, $!";
      next;
      }
   my $lcount = scalar @lines;

   my $lat;
   for(my $i=0; $i<$lcount; $i++) {
      $lines[$i] =~ s/^\s+|\s+$//g;
      my $nb_entries = count_tabs(\$lines[$i]);
      #print "$lines[$i] # $nb_entries\n";
      if ($lines[$i] =~ m/^$/) { # empty line
         push @$lats_ref, $lat; # store lattice hash in the array
         }
      elsif ($nb_entries == 0 && $i+1 < $lcount && $lines[$i+1] !~ m/^$/) { # header
         $lines[$i] =~ m/---([\.0-9]+)-/;
         if (! defined ($1)) {
            print STDERR "$0: cannot detect start time in the file name '$lines[$i]'\n";
            exit -1;
            }
         $lat = LatGraph->new($lines[$i], $frameRate, $1, $frameBased);
         }
      elsif ( $nb_entries == 3 ) { # edge / word     
         if ($frameBased == 1) {   
            my $nb_frames = count_underscores(\$lines[$i]) + 1;
            $lines[$i] =~ m/^(\d+)\t(\d+)\t([[:alnum:]<>]+)/;
            $lat->add_edge($1, $2, $3, $nb_frames); 
            }
         elsif ($frameBased == 0) {   
            $lines[$i] =~ m/^(?:\d+)\t(?:\d+)\t(?:\d+)\t(?:[-e\.0-9]+),(?:[-e\.0-9]+),([_0-9]+)/; 
            if (!defined $1) { 
            print "$lines[$i] -- $1\n";
            }
            # $1 -> 419_493_565_369_620
            $lat->add_phones($1); 
            }
         }
      }
   }
#-----------------------------

#sub sort_by_time {
#   my ( $lats_ref ) = @_;

#   sort { $a->{_startTime} <=> $b->{_startTime} } @{$lats_ref};

#   }

sub label_best {
   my ( $lats_ref, $lats_best_ref ) = @_;

   for (my $i=0; $i<scalar @$lats_ref; $i++) { 
      $$lats_ref[$i]->label_best($$lats_best_ref[$i]);
      }
   }
#-----------------------------

sub label_phones {
   my ( $lats_ref, $lats_phones_ref ) = @_;

   for (my $i=0; $i<scalar @$lats_ref; $i++) { 
      if ((scalar @{$$lats_ref[$i]->{_e}}) != (scalar @{$$lats_phones_ref[$i]->{_e}})) {
         print STDERR "$0: Different number of lattice L1/L3 edges in '$$lats_ref[$i]->{_name}'\n";
         exit -1;
         }
      $$lats_ref[$i]->label_phones($$lats_phones_ref[$i]);
      }
   }
#-----------------------------

sub reduce_out {
   my ( $lats_ref ) = @_;

   for (my $i=0; $i<scalar @$lats_ref; $i++) { 
      my $keys_ref = $$lats_ref[$i]->reduce_out();
      $$lats_ref[$i]->{_name} =~ m/-S([0-9]+)---/;
      if (! defined ($1)) {
         print STDERR "$0: cannot detect speaker name in the file name '$$lats_ref[$i]->{_name}'\n";
         exit -1;
         }
      print '# '.($i+1).' S'.$1."\n";
      $$lats_ref[$i]->print_sorted($keys_ref, 'STO');
      }
   }
#-----------------------------

sub print_lats {
   my ( $lats_ref, $options ) = @_;
   for (my $i=0; $i<scalar @$lats_ref; $i++) { 
      $$lats_ref[$i]->print($options);
      }
   }
#-----------------------------

sub read_map {
   my ( $filename, $map_ref ) = @_;

   # Read word-map
   my $F;
   my $rc = (open $F, '<:utf8', $filename) || die "Couldn't open file $filename, $!"; 
   while(<$F>) {
      my @A = split(" ", $_);
      @A == 2 || die "bad line in symbol table file: $_";
      $$map_ref{$A[1]} = $A[0];
      }
   close($F);
}
#-----------------------------

read_map($main::symtab_w, \%main::i2w);
read_map($main::symtab_p, \%main::i2p);

my @L1;
my @L2;
my @L3;

read_lats($L1, \@L1, 0.03, 1); # frameBased
#print_lats(\@L1, 'TW');
#print Dumper \@L1;
read_lats($L2, \@L2, 0.03, 1); # frameBased
read_lats($L3, \@L3, 0.03, 0); # phoneBased
#print_lats(\@L2, 'TW');

# Check if the number of items in both lattice collections is the same
if (scalar @L1 != scalar @L2) {
   print STDERR "Error: Lattice collections $L1 and $L2 have different sizes\n";
   exit -1;
   }
if (scalar @L1 != scalar @L3) {
   print STDERR "Error: Lattice collections $L1 and $L3 have different sizes\n";
   exit -1;
   }

# check if filenames are identical and in the same order in both collections
for(my $i=0; $i<scalar @L1; $i++) {
   if ($L1[$i]->{_name} ne $L2[$i]->{_name}) {
      print STDERR "Error: Lattice collections $L1 and $L2 have mismatching names '$L1[$i]->{_name}' and '$L2[$i]->{_name}'\n";
      exit -1;
      }
   }
for(my $i=0; $i<scalar @L1; $i++) {
   if ($L1[$i]->{_name} ne $L3[$i]->{_name}) {
      print STDERR "Error: Lattice collections $L1 and $L3 have mismatching names '$L1[$i]->{_name}' and '$L3[$i]->{_name}'\n";
      exit -1;
      }
   }

# mark best hypothesis in the first colection of lattices
label_best(\@L1, \@L2);
label_phones(\@L1, \@L3);
#sort @L1;
my @s_L1 = sort { $a->{_startTime} <=> $b->{_startTime} } @L1;
reduce_out(\@s_L1);
#print_lats(\@L1, 'STO');
#print_lats(\@L1, 'SNTIWP');

exit 0;


# nereikalinga, jei naudojame lattice-1best, o ne lattice-to-nbest --n=1
# strip '-1' from the ending of the filename in the 1-best collection
# for(my $i=0; $i<scalar @latnames2; $i++) {
#   $latnames2[$i] =~ s/-1$//;
#   }


while (<>) {
   my @A = split(" ", $_);
   for (my $pos = 0; $pos <= $#A; $pos++) {
      my $a = $A[$pos];
      if ($pos == 2) {
         $a = int2word($a);
         }
      elsif ($pos == 3) {
         my @B = split(",", $A[$pos]);
         $a = $B[0].",".$B[1];
         my @C = split("_", $B[2]);
         my $transc='';
         for (my $p = 0; $p <= $#C; $p++) {
            $transc = $transc.int2phone($C[$p], 1); # remove_pd_info
            }
         $a = $a." ".$transc;
         }
      print "$a ";
      }
   print "\n";
   }

exit 0;

