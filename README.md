ConcertWatch
===
Searching through the ticketmaster website manually is an absolute pain, so instead
this script will query the ticketmaster API for a list of nearby concerts and then
compare the artists with a list of artists I want to be notified about.

Long-term I would like to add some of the following features:
- Generate the list of artists from youtube channels an acount is subscribed to.
- Generate the list of artists based on songs that a Youtube channel has liked.
- Diff mode that only sends notifications based on events we haven't yet seen.
- Search Criteria customisation via YAML file instead of code.
