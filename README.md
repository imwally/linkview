# linkview

Like [urlview](https://github.com/sigpipe/urlview) but tailored for HTML documents.

That is, unlike `urlview` it only reads valid HTML documents as it
relies on the [goquery](https://github.com/puerkitobio/goquery)
package to parse the HTML and find `a` elements. The advantage to this
is being able to extract the text from a link or in the case of an
image, the alt or title attributes. This gives more context to what
the URL points to. Link text will show in the menu instead of just the
URL.

Here's an example:

```
j/C-n: move down   k/C-p: move up   return/C-o: open url   q: quit

https://t.e2ma.net/click/w74ifb/cow3ngb/023tkn

   Longwood Gardens
   social-email.png
   social-twitter.png 
   social-facebook.png
   daily
-> Buy Tickets
   ftr_01.png
   1001 Longwood Road, Kennett Square, PA 19348
   longwoodgardens.org
   Facebook
   Twitter
   Instagram
   YouTube
   Blog
   Manage
   Opt out
   Sign up
   online
```

The first section displays help text, the second displays a URL
preview of the currently selected link, and finally the menu of
links. Pressing the `return` key will open the currently selected link
in your default browser.

## But Why?

Many HTML emails are built using platforms that modify original
URLS. This makes it hard to figure out what the URL points to. If you
notice in the example above, the URL to Buy Tickets is obscured by
some marketing URL redirect.
