# EasyRM

EasyRM is a personal relationships manager that sits on top of Google Contacts (for now).

It is meant to help you keep in touch with friends, and to keep all your friends' information in a central
and accessible location. I chose Google Contacts as a backend because it already syncs across devices, and it
integrates with other Google services in useful ways (e.g. a contact's address will show up on Google Maps,
and their birthday on Google Calendar).


## Merging Facebook friends into EasyRM

One key feature of a relationship manager is to unify a contact's information in one place. To that end, I
want my EasyRM contacts to link to their other social media profiles. Here's how to get your friends' Facebook
URLs and merge them into your EasyRM contact.

1. Go to your friends page: `https://facebook.com/<your-username>/friends`

2. Right-click on a friend's name, and click on "Inspect". On the HTML view, right-click and select "Copy full Xpath".

The xpath will look something like this:
```
/html/body/div[1]/div/div[1]/div/div[3]/div/div/div/div[1]/div[1]/div/div/div[4]/div/div/div/div/div/div/div/div/div[3]/div[12]/div[2]/div[1]/a/span
```

At the time of writing, the HTML was structured like so:
```html
<a href="https://facebook.com/<username>" ...>
    <span ...>Friend's Name</span>
</a>
```

Thus, we'll make the following modifications to the XPATH (remote the last `span` and match all friends by removing `[12]`):
```
/html/body/div[1]/div/div[1]/div/div[3]/div/div/div/div[1]/div[1]/div/div/div[4]/div/div/div/div/div/div/div/div/div[3]/div/div[2]/div[1]/a
```

3. On the console, run the following:

```javascript
links = $x("/html/body/div[1]/div/div[1]/div/div[3]/div/div/div/div[1]/div[1]/div/div/div[4]/div/div/div/div/div/div/div/div/div[3]/div/div[2]/div[1]/a");
pairs = links.map((anchor) => [anchor.href, anchor.childNodes[0].innerHTML]);
```

Now, `pairs` contains a list of lists (or tuples), where each tuple is (facebook profile url, name).

4. Get all friends

Now that we know we can get the url/name from the friends page, we need to make sure that the page
actually contains all of our friends. The page uses infinite scroll, and the only way I've found
how to do this, is to manually keep scrolling down until it ends. It will take a few minutes of
mindless scrolling til you reach the end.

Once you do, run the code above again and you'll have a list of url/names for all your friends!

5. Save it. Run `people.MergeFacebookURLs` (to be moved to some sort of CLI at some point). Voila!

### Inactive profiles

Some of your friends may be inactive. You can use this code to find them:
```javascript
all = $x("/html/body/div[1]/div[1]/div[1]/div/div[3]/div/div/div/div[1]/div[1]/div/div/div[4]/div/div/div/div[1]/div/div/div/div/div[3]/div");
all.forEach(div => {
  if (div.childNodes[1].childNodes[0].childNodes[0].tagName !== 'A') {
    console.log(div.childNodes[1].childNodes[0].childNodes[0].innerHTML);
  }
});
```

