#!/usr/bin/perl

#XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
#X                                                                  X
#X lat_restore.pl, Copyright(C) Gailius Raðkinis, 2020-2022         X
#X                                                                  X
#XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX

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
use Digest::MD5  qw(md5_hex);
use lib dirname (__FILE__); # kad rastu .pm
use LatGraph;
#use Data::Dumper qw(Dumper);

my $dirname = dirname(__FILE__);
require "$dirname/lat_map.pl";
require "$dirname/utils_num.pl";
require "$dirname/ph2word.pl"; 

binmode(STDIN, ":utf8");
binmode(STDOUT, ":utf8");

# TODO kai --join-num 0.1 $deb_lat = 139 duoda "Lattice no. 139 has 1 broken segments"

our $debug = 0;
my $deb_lat; # = 139;# = 96;# = 45;# = 184; # = 190;# = 183; # = 192; # = 33; ## = 210; # = 39; #203;#  = 190;# = 208;
# 33 - <eps>
# 39 - 600 vs. 649
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

# joins two adjacent segments if they belong to the same speaker and 
# their summary length is inferior to $max_sum_duration sec (long text is bad for VVT)
our $join_spk = 0;
our $max_sum_duration = 100000.0; # default join all sections belonging to the same speaker 
# joins adjacent numbers e.g. 9 1000 6 100 50 2 -> 9652
# if they are not separated by silence longer than $max_intra_sil sec
our $join_num = 0;
our $max_intra_sil = 0.03;
# Semantikos-2 redaktorius rodo hipoteziu alternatyvas
# Intelektikos-2 redaktorius alternatyviu hipoteziu nerodo
# Taupant skaiciavimo laika (FST <unk> apdorojimas suletino ðá komponeta) ir
# siekiant minimizuoti <unk> zodziu kieki, keiciama elgsena pagal nutylëjimà
# Nuo siol alternatyvios hipotazes eliminuojmos, nebent nurodoma '--keep-alt'
our $keep_alt = 0;

while (scalar @ARGV > 0) {
   if ($ARGV[0] eq '--join-spk') {
      $join_spk = 1;
      shift @ARGV;
      if (defined $ARGV[0] && $ARGV[0] =~ m/^[0-9\.]+$/) {
         $max_sum_duration = $ARGV[0]; 
         shift @ARGV;
         }
      }
   elsif ($ARGV[0] eq '--join-num') {
      $join_num = 1;
      shift @ARGV;
      if (defined $ARGV[0] && $ARGV[0] =~ m/^[0-9\.]+$/) {
         $max_intra_sil = $ARGV[0]; 
         shift @ARGV;
         }
      }
   elsif ($ARGV[0] eq '--keep-alt') {
      $keep_alt = 1;
      shift @ARGV;
      }
   else { # not accepted
      shift @ARGV;
      }
   }

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
         $lines[$i] =~ m/---([\.0-9]+)-([\.0-9]+)/;
         if (! defined ($1) || ! defined ($2)) {
            print STDERR "$0: cannot detect start/end times in the file name '$lines[$i]'\n";
            exit -1;
            }
         $lat = LatGraph->new($lines[$i], $frameRate, $1, $2, $frameBased);
         }
      elsif ( $nb_entries == 3 ) { # edge / word     
         if ($frameBased == 1) {   
            my $nb_frames = count_underscores(\$lines[$i]) + 1;
            $lines[$i] =~ m/^(\d+)\t(\d+)\t([[:alnum:]<>]+)\t([-e\.0-9]+,[-e\.0-9]+,[_0-9]+)$/;
            # $4 - footprint must include not only transition ids but probs as well (identical <unk> transition ids observed)
            $lat->add_edge($1, $2, $3, md5_hex($4), $nb_frames); 
            }
         elsif ($frameBased == 0) {   
            $lines[$i] =~ m/^(?:\d+)\t(?:\d+)\t(?:\d+)\t(?:[-e\.0-9]+),(?:[-e\.0-9]+),([_0-9]+)/; 
            if (!defined $1) { 
               print STDERR "$0: Unable to extract phone ids from line ".$i." '$lines[$i]'\n";
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
      next if (defined $deb_lat && $i != $deb_lat); # debug
      $$lats_ref[$i]->label_best_v3($$lats_best_ref[$i]);
      $$lats_ref[$i]->unlink_alternatives() if ($keep_alt == 0);
      }
   }
#-----------------------------

sub label_phones {
   my ( $lats_ref, $lats_phones_ref ) = @_;

   for (my $i=0; $i<scalar @$lats_ref; $i++) { 
      next if (defined $deb_lat && $i != $deb_lat); # debug
      if ((scalar @{$$lats_ref[$i]->{_e}}) != (scalar @{$$lats_phones_ref[$i]->{_e}})) {
         print STDERR "$0: Different number of lattice L1/L3 edges in '$$lats_ref[$i]->{_name}'\n";
         exit -1;
         }
      $$lats_ref[$i]->label_phones($$lats_phones_ref[$i]);
      }
   }
#-----------------------------

sub connect_nums {
   my ( $lats_ref ) = @_;

   for (my $i=0; $i<scalar @$lats_ref; $i++) { 
      next if (defined $deb_lat && $i != $deb_lat); # debug
      $$lats_ref[$i]->connect_nums();
      }
   }
#-----------------------------

sub reduce_out {
   my ( $lats_ref ) = @_;

   my $last_spk = -1;
   my $lasti    = -1;
   for (my $i=0; $i<scalar @$lats_ref; $i++) { 
      next if (defined $deb_lat && $i != $deb_lat); # debug
      next if ($$lats_ref[$i]->is_empty() == 1); # skip initial empty lattices
      my $keys_ref = $$lats_ref[$i]->reduce_out();
      $$lats_ref[$i]->{_name} =~ m/-S([0-9]+)---/;
      if (! defined ($1)) {
         print STDERR "$0: cannot detect speaker name in the file name '$$lats_ref[$i]->{_name}'\n";
         exit -1;
         }
      if ($join_spk == 0 || $last_spk == -1 || $1 != $last_spk || $$lats_ref[$i]->{_endTime} - $$lats_ref[$lasti]->{_startTime} > $max_sum_duration) {
         print "\n" if ($i > 0);         # newline between sections
         print '# '.($i+1).' S'.$1."\n"; # new section header
         $lasti = $i;
         }
      $last_spk = $1;
      $$lats_ref[$i]->print_sorted($keys_ref, 'STO');
 
      # print following empty lattices
      my $j;
      for($j=$i+1; $j<scalar @$lats_ref; $j++) {
         last if ($$lats_ref[$j]->is_empty() == 0); 
         $keys_ref = $$lats_ref[$j]->reduce_out();
         $$lats_ref[$j]->print_sorted($keys_ref, 'STO'); 
         }
      $i = $j-1;
      }
   }
#-----------------------------

sub print_lats {
   my ( $lats_ref, $options ) = @_;
   for (my $i=0; $i<scalar @$lats_ref; $i++) { 
      next if (defined $deb_lat && $i != $deb_lat); # debug
      $$lats_ref[$i]->print($options);
      }
   }
#-----------------------------

# Perform checks to see if modified lattices satisfy necessary conditions
sub test_lats {
   my ( $lats_ref ) = @_;
   for (my $i=0; $i<scalar @$lats_ref; $i++) { 
      next if (defined $deb_lat && $i != $deb_lat); # debug
      $$lats_ref[$i]->test_single_best($i);
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
my @s_L1 = sort { $a->{_startTime} <=> $b->{_startTime} } @L1;
# Connect digits
print_lats(\@s_L1, 'FSNTIWP') if ($debug == 1);

connect_nums(\@s_L1) if ($join_num == 1);
# Print text hypotheses
reduce_out(\@s_L1); # cia spausdina rezultata
#print_lats(\@L1, 'STO');
#print_lats(\@L1, 'SNTIWP');

test_lats(\@s_L1);

exit 0;

