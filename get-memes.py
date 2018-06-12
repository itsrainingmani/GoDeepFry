import praw
import pprint
import requests
import re

# function that checks the content type from the headers and determines if it is downloadable or not
def is_downloadable(url):
    h = requests.head(url, allow_redirects=True)
    header = h.headers
    content_type = header.get('content-type')
    if 'text' in content_type.lower() or 'html' in content_type.lower():
        return False
    return True

# function that gets the filename from the url. The filename is the characters to the right of the last backslash
def get_filename_from_url(url):
    if url.find('/'):
        return url.rsplit('/', 1)[1]

# Initialize the reddit praw instance with the credentials defined in the praw.ini files
reddit = praw.Reddit('GoDeepFry', user_agent='windows:GoDeepFry:v0.1 (by /u/L-king)')

subreddit = reddit.subreddit('dankmemes')
folder = "./memes/"
for submission in subreddit.top(limit=100):
    # We only want the submissions that aren't v.redd.it videos, gifs or gifvs
    if submission.is_video == False and 'gifv' not in submission.url and 'gif' not in submission.url:
        if is_downloadable(submission.url):
            r = requests.get(submission.url, allow_redirects=True)
            filename = folder + get_filename_from_url(submission.url)
            open(filename, 'wb').write(r.content)