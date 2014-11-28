gotags
======

gotags is a ctags replacement. The newer ctags does support Go but it has some
limitations. For example, it will not tag struct field names, but it will tag
the struct name. There are other projects that add editor support for Go, e.g.,
supporting Omni-Complete in Vim, but I only wanted ctags.

So, if you just need ctags for Go you can build then run this tool. It is
recursive by default so no need to pass a -R option.
