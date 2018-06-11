import praw

reddit = praw.Reddit('GoDeepFry', user_agent='windows:GoDeepFry:v0.1 (by /u/L-king)')

subreddit = reddit.subreddit('dankmemes')
for submission in subreddit.top(limit=10):
    print(submission.title)
    print(submission.score)