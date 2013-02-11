Neobin
========

Neobin is a very simple but efficient binary format.

It was created to fill the gaps left behind by binary formats like protobuf, that were not designed to handle huge files, or files with millions of small data (two floats).

Neobin format is very efficient at reading huge data files, for instance matrix format.

Neobin binary files always starts with a **magic number** to identify "neobin", then a **unique namespace** string declared in the file format specification. Neobin binary file can be introspected to find their format (just like xml namespaces).


The generator comes as a simple executable

The fileformat Schema (XSD is available <a href="http://code.google.com/p/neobin/source/browse/src/main/resources/neobin_v1.xsd"> here</a> ) gets you  autocompletion in your favorite IDE

Both the fileformat, and the generated binary files can be read to and from <a href="http://code.google.com/p/neobin/">java</a>.

Visit the <a href="http://code.google.com/p/neobin/">neobin-java</a> website to learn more about this format.


The code is available at http://www.gopack.me
<pre>
git clone git://github.com/ericaro/goneobin.git
cd goneobin
gpk compile
</pre>

