# Project Spring-bonfire 0219309 AM/PM AkA Project Phoenix AkA TagMachine

![Untitled drawing (7)](https://github.com/hartsfield/machineTag/assets/30379836/87bfbd09-ed33-4584-8ca9-d1a7034ee9ab)

# https://TagMachine.xyz/

minimum requirements:


    # Install these (using homebrew on a mac)
    git   # the git program
    gh    # the github command line app
    redis # database software for tagmachine
    
    # Install Go
    https://go.dev/doc/install
    
    # Download and run tagmachine
    cd                                                        
    git clone https://github.com/hartsfield/tagmachine.xyz
    cd tagmachine.xyz
    ./autoload.sh tagmachine.xyz 9002 SECRET_PHRASE_HERE
    
    # you should be presented with the address the app is running on like:
    > > > > http://localhost:9002/
    
    # Start Redis
    idk how on a mac actually, it usually just starts itself

    # Hierarchic explanations
    <tagmachine.xyz/

          ▾ internal/         # The files contained in internal/* are "templates", which could be either "components" or "pages"
            ▸ components/     # where the components live
            ▸ pages/          # where the pages live
            ▸ shared/         # where we keeps shared resources, which can be accessed by all pages and components
                              # NOTE: Components can be made up of smaller components, or added to a page, but a page is never "added" to
                              # a component, or another page

          ▸ public/           # public files, where we store/save media

          authHelpers.go      # the rest should be *fairly* self explanatory
          checkAuth.go        # FYI: 1075 lines of Go code, and about as much 
          globals.go          # html/css/js in the templates.
          handlers.go         # The entire app is only about 2000 lines of sourcecode (maybe)
          helpers.go
          main.go
          router.go
          server.go
          signin.go
          signout.go
          signup.go
          viewdata.go
          autoload.sh
          bolt.conf.json
          LICENSE
          README.md
          tagmachine

# This is a re-write of TagMachine (https://github.com/hartsfield/tagmachine)

NOTE: This program is still in alpha phases and is extremely unstable. Don't 
expect it to be bug free. 

TagMachine is a new type of social media website that aims to accept and 
embrace social media as a modern form of journalism. 

## TagMachine Philosophy
TagMachine is designed to be a hybrid, impersistent, social media platform for people seeking discussion of a broad set of topics, including up-to-the-minute breaking news and world events.

We believe social media can be used as a modern form of decentralized journalism, and aim to create a platform to quickly relay and funnel real-time world events into an easily digested stream of information. TagMachines goal is to replace "mass-media" with something better, more resilient, and more decentralized.

New generations are increasingly consuming news via social media, but most social media (facebook, tiktok, etc) are failing us as gatekeepers. TagMachine views itself as a type of modern journalism platform for all topics. As gatekeeper to this platform we believe it's important to enforce a high standard for discussion.

TagMachine has a similar philosophy to Ycombinators Hacker News, but we'd like to cater to a wider variety of topics and people. We want to foster serious discussion and activism, and eventually expand the platform to increase it's relevancy as a replacement for mass-media and traditional news sources. This includes encourging use of our platform in places around the world, including Asia, Africa, and South America.

Post moderation is not strict. Pretty much anything goes as long as it's on topic and doesn't violate the rules (see below). TagMachine has ABSOLUTELY NO political affiliation and does not endorse ANY political party or stance, however we strongly encourge healthy debates and arguments on our platform, and we believe our platform is designed to foster this type of discussion.

Users who participate in debates in an un-civil manner will risk having their accounts suspended or permanently deleted. A good argument will provide facts and sources, and should not resort to ad-hominem.

If TagMachine's algorthms work properly, it should create healthy and active discussions.

## RULES
No spam/advertising
No bots or organized information suppression
No sexualization of minors
No racism/sexism
No links to illegal content
No posting personal information of anyone without consent
No harassment of other users

## TAGS, TOPICS, SUBREDDITS
Much like twitter, users have a @handle and can #tag their posts with different topics. In TagMachine, these topics have a popularity score and their popularity is based on their frequency of use. Unlike twitter, TagMachine has a set of default tags that we believe cover almost all topics, and these tags are persistent.

Users are required to pick at least one of these tags, and provide at least one non-default tag when creating a new thread. The user-generated tags, along with the default tags, are displayed on the "frontpage" of TagMachine, and will eventually be used in a word cloud to visualize trends and breaking news happening around the world. Once tuned correctly, this feature of TagMachine will be able to report on breaking news at a faster rate than traditional news platforms.

Only tags with a score above a certain threshold are displayed in the tag cloud.

Tags are used to filter topics and are somewhat analogous to reddits subreddits.

## SCORES, LIKES, UPVOTES
Everything has a score.

Tags have a score, users have a score, and posts have a score.

A score is analagous to reddit karma.

Unlike most social media sites, a posts score is only increased when another user replies to it. The score also "bubbles up" the comment tree, meaning that when a user replies to a comment, that comment and the person who posted it receive a score increase, and so do all parent comments all the way up to the original post and the tags posted with it.

There is no way to "downvote" or "dislike" a thread. If you don't like a topic, either don't reply, or reply with a well thought out argument against the topic. If users in the thread disagree with you, all of their replies will increase your comments score (and your user score). This could have a chain-reaction, pushing your comment to the top and gaining it even more replies. Users who engage in dialogue are rewarded by TagMachines algorithm. On reddit you would likely be downvoted into oblivion. On most other platforms only the most agreeable comments will be seen. TagMachine is designed so that the most active conversations rise to the top, even when people disagree.

## IMPERSISTENCE
Like many anonymous image boards, TagMachine is not persistent. This means that threads are deleted automatically from the database. Most threads should not last more than 72 hours. Points accumulated by users and tags are kept. User-generated tags have their scores decreased over time if users don't continue submitting the tag. Unlike most anonymous image boards, we currently require an account to post, and we enforce more strict moderation guidelines.

## NO EMBEDDED MEDIA
As of version 1.0 there is no embedded media allowed. This is for a few reasons:

Website performance and costs. Allowing embedded media would cost exponentially more.
We aren't trying to be another tiktok, in fact, we are trying to make a conscious effort to be something different.

Links are likely to be filtered soon too. Instead, we'll have a whitelist of whats considered a good source (eg scholarly publications)

## SOFTWARE
TagMachine is not a clone of some other social media platform, but rather a distilled "hybrid" platform, taking some features and concepts from all of the most popular social media platforms, and combining them into a new concept, that we believe is more efficient, and allows for greater insight and transparency into the community.

TagMachine follows KISS and DRY principals of software development. It is mostly written from scratch using the Go programming language, and uses component architecture via Go HTML templating. The code for TagMachine will eventually be released to the public, but for now it's only partially available upon request.
 
## Features implemented so far:
 - User Authentication
 - Posts
 - Replies
 - Webm/mp4/jpeg/png/gif upload support
 - Ranked posts
 - Two sorted views (ranked/new)
 - Pagification

SLOC = 1601 as of this commit

## Features implemented in the old version that need to be ported:
 - Tagging
 - Filter by tag
 - Mentions
 - Follow user
 - User pages
     - following stream
     - suggestions

## Features coming soon:
 - User favorites
 - Maybe votes
 - Maybe mp3 support
 - Maybe sorting algorithms

# Development
## Requirements
  - Go programming environment
  - Redis (v7+)
  - Only tested on Linux (Debian & RedHat based)

## Instructions for running

Clone the repo, and run the following (assuming redis is running on the default port):

    go mod tidy
    chmod +x autoload.sh
    ./autoload.sh tm2 9999 [random phrase here]

Now visit `http://localhost:9999`, sign-up, and add some posts.

## Bolt Architecture

An explanation of what Bolt architecture is can be found in [the git repository for the Bolt software](https://github.com/hartsfield/bolt)

Basically the front end is built using components composed of regular 
HTML/css/JavaScript, but they're executed via Go html templating, so Go 
template directives can be used. 

 - These components are mostly self-contained in directories located in `internal/components`.
 - Once created, these components can be added to a `page` located in 
`internal/pages`, or they can be used as a sub-component in another component.
 - Routes are registered in `routes.go`
 - View/Model structs are located in `viewdata.go`

![Screenshot from 2023-08-26 19-23-38](https://github.com/hartsfield/machineTag/assets/30379836/6fa734ad-2dfb-4387-8f24-d8386acec19c)
