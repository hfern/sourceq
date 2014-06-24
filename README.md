sourceq
=======

Command line query tool for source servers. Basically a head for (goseq)[https://github.com/hfern/sourceq].
Use it to get source server lists from the Valve's master servers (or
a custom, ip-specified one).

Installation
============

Usage
=====

    #General help for querying master servers
    sourceq master -?
    
    # Get the first 20 servers' IP and Name in U.S. West.
    sourceq master -fields "ip,name" -l20 -r"USW" -a
    
    # List the usable filters, regions, and fields that can be used.
    sourceq master --list-filters --list-regions --list-fields

