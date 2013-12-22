Flower
======

The Flower extensions to Blackfriday, bring executable network specifications
to [Markdown][1] documents. 

These extenisons are only enabled if the `--flower=true` command line option is specified. 


Design
------

Code blocks are passed through a wrapper that parses the text looking for the syntax `flower: …` which signifies a flower extension.

The flower code is surrounded by a `<div id='flower-XX'>` where XX is a unique number for each flower command in the document.  

Commands are processed in the order they are seen.

Each command is evaluated to see if it applies to the current host, either as the host offering a service, or consuming a service. If it doesnt apply the command has a result NO_ANSWER.

If the command does apply, then the network connection is tested to see if it works, by scanning to see if that port is open and accessable, and will have a result of OK, FAIL or ERROR.

The flower code `<div>`'s are styled using CSS accoring to their sucess or failure.
