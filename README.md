threepio
==

![threepio](https://s-media-cache-ak0.pinimg.com/736x/c5/35/c9/c535c913ca0135bd19010f013a7e65f6.jpg)

threepio is a little golang app. It responds to a URI string dictating where an asset for editting is stored, syncs this and opens either prelude or premiere for editting.

Grammar
--

```
uri             = "threepio+", application, "://", path, "?", params;
application     = "prelude" | "premiere";

path            = "/", alphanumeric, {path};

params          = param, {"&", params};

param           = key, "=", value;
key             = "accessKey" | "secretKey" | "sessionToken" | "uuid";
value           = alphanumeric;


alphanumeric    = letter | digit, {alphanumeric | punctuation};
```

Building
--

To test/ run the threepio binary:

```bash
$ go build threepio.go
```

To build the actual `.dmg` installer, and register the protocol handler, the procedure is a bit more involved.

Ultimately you'll need to follow [this blog post](http://ramezanpour.net/post/2014/05/12/how-to-create-installation-dmg-files-in-os-x/)
