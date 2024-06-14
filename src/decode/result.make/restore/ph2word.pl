#!/usr/bin/perl

use strict;
use warnings;
#use Encode qw(decode encode);
use LWP::UserAgent;
use JSON;
use JSON::Parse 'parse_json'; 

# https://www.xmodulo.com/how-to-send-http-get-or-post-request-in-perl.html

my $ua = LWP::UserAgent->new;
my $server_endpoint = "http://127.0.0.1:3000/phones2word";

# set custom HTTP request header fields
my $req_get = HTTP::Request->new(GET => $server_endpoint);
$req_get->header('content-type' => 'application/json; charset=utf-8');

my $req_post = HTTP::Request->new(POST => $server_endpoint);
$req_post->header('content-type' => 'application/json; charset=utf-8');

sub phones2letters {
   my ( $phstr ) = @_;

   my %data = ( 'phones' => $phstr );
   $req_get->content(encode_json(\%data));

   my $resp = $ua->request($req_get);
   if ($resp->is_success) {
      my $json = $resp->decoded_content;
      my $perl = parse_json($json);
      return decode('utf-8', $perl->{word});
      }
   else {
      return(p2g_rules($phstr));
      }
   }
#-----------------------------

sub phones2letters_list {
   my ( $list_ref ) = @_;

   $req_post->content(encode_json($list_ref));

   my $resp = $ua->request($req_post);
   if ($resp->is_success) {
      my $content = $resp->decoded_content;
      my $json = parse_json($content);
      for(my $i=0; $i<scalar @$json; $i++) {
         $list_ref->[$i]->{word} = decode('utf-8', $json->[$i]->{word});
         # print $list_ref->[$i]->{word}."\n"; # debug
         }
      }
   else {
      foreach (@$list_ref) {
         $_->{word} = p2g_rules($_->{phones})        
         }
      }
   return($list_ref);
   }
#-----------------------------

sub p2g_rules {
   my ( $phones ) = @_;

   no warnings qw(once);
   $phones =~ s/(?<= [aeo] v\' i tS\' )\"e:(?= m. s$)/ia/g; # stankeviciams
   $phones =~ s/(?<= ^[aeo]: v\' i tS\' )eu/iau/g; # stankeviciaus

   $phones =~ s/tS/$::ch/g;
   $phones =~ s/Z/$::zh/g;
   $phones =~ s/S/$::sh/g;
   $phones =~ s/E/$::eh/g;
   $phones =~ s/ts/c/g;
   $phones =~ s/G/h/g;
   $phones =~ s/x/ch/g;
   $phones =~ s/N/n/g;

   $phones =~ s/(?<=j)( [\^\"]?)i([ou])/$1$2/g;

   $phones =~ s/(?<= )a: j i:$/$::a_nosine j $::i_nosine/g; # aji
   $phones =~ s/(?<= )a: j e:$/$::a_nosine j $::a_nosine/g; # aja
   $phones =~ s/(?<= )i: j i:$/$::i_nosine j $::i_nosine/g; # iji
   $phones =~ s/(?<= )u: j i:$/$::u_nosine j $::i_nosine/g; # uji
   $phones =~ s/( [\^\"]?i?)u: j u:$/$1$::u_nosine j $::u_nosine/g;# uju

   $phones =~ s/[\^\"]i:$/y/g;
   $phones =~ s/(?<= )i:$/$::i_nosine/g;
   $phones =~ s/(^|^[nbt]\' e |^[nt]\' e b\' e)[\^\"]?i:/$1$::i_nosine/g;
   $phones =~ s/i:/y/g;

   $phones =~ s/(?<= )[\^]?a:$/$::a_nosine/g;
   $phones =~ s/(?<= )[\^]?a:(?= s$| s\' i$)/$::a_nosine/g;
   $phones =~ s/(?<= )e:(?=$| s$| s\' i$| s\' i s$)/$::e_nosine/g;

   $phones =~ s/( [\^]?i?)u:$/$1$::u_nosine/g;
   $phones =~ s/((?:^| )[\^\"]?i?)u:(?= )/$1$::u_ilgoji/g;

   $phones =~ s/[\.:\'\"^ ]//g; # delete single quote

   return $phones;
   }
#-----------------------------

sub is_unk {
   my ( $word ) = @_;

   return 1 if ($word eq '<unk>' || $word =~ m/^_.+_$/);
   return 0;
   }
#-----------------------------

1

