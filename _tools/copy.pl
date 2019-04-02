#!/usr/bin/env perl
use 5.014;
use warnings;
use utf8;
use autodie;
use File::Spec;
use File::Copy qw/copy/;
use File::Path qw/mkpath/;
use File::Basename qw/basename/;

my $goroot = 'github.com/golang/go/src';
chomp(my $ghq_root = `ghq root`);
my $modfile_dir = File::Spec->catfile($ghq_root, $goroot, 'cmd/go/internal/modfile');

opendir my $dh, $modfile_dir or die $!;
while (my $f = readdir $dh) {
    next if $f !~ /\.go$/;
    next if $f =~ /_test\.go$/;
    copy(File::Spec->catfile($modfile_dir, $f), '.');
}
closedir $dh;

for my $dir (qw!cmd/go/internal/semver cmd/go/internal/module internal/lazyregexp!) {
    mkpath basename($dir);
    my $pkg_dir = File::Spec->catfile($ghq_root, $goroot, $dir);
    opendir my $dh, $pkg_dir or die $!;
    while (my $f = readdir $dh) {
        next if $f !~ /\.go$/;
        next if $f =~ /_test\.go$/;
        copy(File::Spec->catfile($pkg_dir, $f), "$dir/");
    }
    closedir $dh;
}
