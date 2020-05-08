#!/usr/bin/perl

use strict;
use warnings;
use Encode qw(decode encode);

our $a_nosine = decode("utf-8", "\xC4\x85");
our $ch       = decode("utf-8", "\xC4\x8D");
our $e_nosine = decode("utf-8", "\xC4\x99");
our $eh       = decode("utf-8", "\xC4\x97");
our $i_nosine = decode("utf-8", "\xC4\xAF");
our $sh       = decode("utf-8", "\xC5\xA1");
our $u_nosine = decode("utf-8", "\xC5\xB3");
our $u_ilgoji = decode("utf-8", "\xC5\xAB");
our $zh       = decode("utf-8", "\xC5\xBE");

our $A_nosine = decode("utf-8", "\xC4\x84");
our $CH       = decode("utf-8", "\xC4\x8C");
our $E_nosine = decode("utf-8", "\xC4\x98");
our $EH       = decode("utf-8", "\xC4\x96");
our $I_nosine = decode("utf-8", "\xC4\xAE");
our $SH       = decode("utf-8", "\xC5\xA0");
our $U_nosine = decode("utf-8", "\xC5\xB2");
our $U_ilgoji = decode("utf-8", "\xC5\xAA");
our $ZH       = decode("utf-8", "\xC5\xBD");

our $lt_str = "$a_nosine$ch$e_nosine$eh$i_nosine$sh$u_nosine$u_ilgoji$zh";
our $LT_str = "$A_nosine$CH$E_nosine$EH$I_nosine$SH$U_nosine$U_ilgoji$ZH";

1

