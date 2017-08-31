# linkview

_NOTE: This is a work in progess. Some things such as resizing your
terminal or rendering too many links may cause issues._

Like [urlview](https://github.com/sigpipe/urlview) but tailored for HTML documents.

`linkview` parses an HTML document for any links and displays a menu
to choose from but instead of displaying only the URL it attempts to
display the link text.

```
j: move down   k: move up   return: open url   q: quit

  http://pages.news.digitalocean.com/n/R0066U4I7V0030h0EIDvFX0

→ Introduction to Object Storage
  API Documentation
  GET STARTED
  tw_blue.png
  soundcloud.png
  Refer a Friend
  NO TEXT
```

The first section displays help text, the second displays a URL
preview of the currently selected link, and finally the menu of
links. Pressing the `return` key will open the currently selected link
in your default browser.



