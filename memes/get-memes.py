import praw
import pprint
import requests
import re

def is_downloadable(url):
    h = requests.head(url, allow_redirects=True)
    header = h.headers
    content_type = header.get('content-type')
    if 'text' in content_type.lower() or 'html' in content_type.lower():
        return False
    return True

def get_filename_from_url(url):
    if url.find('/'):
        return url.rsplit('/', 1)[1]


reddit = praw.Reddit('GoDeepFry', user_agent='windows:GoDeepFry:v0.1 (by /u/L-king)')

subreddit = reddit.subreddit('dankmemes')
for submission in subreddit.top(limit=100):
    # pprint.pprint(vars(submission))
    if submission.is_video == False and 'gifv' not in submission.url and 'gif' not in submission.url:
        if is_downloadable(submission.url):
            r = requests.get(submission.url, allow_redirects=True)
            filename = get_filename_from_url(submission.url)
            open(filename, 'wb').write(r.content)