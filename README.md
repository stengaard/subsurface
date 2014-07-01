subsurface
==========

Did you also forget to set the right date on your dive logger? If so,
this little tool can help you. Simply load up all the dives in
subsurface, save the logbook (and back it up), figure out the time
diff and do us the move tool. Here is what I did:

    $ subsurface -d 78591h0m0s stengaard.ssrf > stengaard1.ssrf

Then load `stengaard1.ssrf` into subsurface and check everything is
there.


Is this of any use to me?
-------------------------
Most likely not.


Hacking
=======

But I need subsurface test xml files:
 - https://github.com/stengaard/divelog
 - https://github.com/torvalds/subsurface/tree/master/tests
