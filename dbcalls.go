// / Provided Under BSD (2 Clause)
//
// Copyright 2025 Johnathan A. Hartsfield
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS “AS IS”
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.
//
// ////////////////////////////////////////////////////////////////////////////
//
// dbcalls.go is where all direct calls to the database (redis) should be
// housed. Functions wishing to apply database procedures outside of this file
// shouldn't make direct calls to the database, but instead should call
// functions located in this file, or if needed create new ones here. This
// keeps all the real database procedures in one place. The following is a
// breakdown of how the database is configured:
//
///////////////////////////////////////////////////////////////////////////////
//
//              TAGSBYSCORE - KEY to ZSET containing reference keys to tags,
//                            ranked by popularity, determined by algorithm.
//
//                    USERS - KEY to ZSET containing reference keys to user
//                            profile data, ranked by user score (for now).
//
//        [loginEmail]:HASH - KEY to VALUE which is the password HASH
//                            associated with the loginEmail.
//
//                   [HASH] - KEY to VALUE which is the authenticated users
//                            user ID, which is then used to look up the users
//                            profile data.
//
//                [user.ID] - KEY to HASH of the associated users data.
//
//   [user.ID]:POSTSINORDER - KEY to ZSET containing reference keys to a
//                            users posts (IDs) in chronological order.
//
//   [user.ID]:POSTSBYSCORE - KEY to ZSET containing reference keys to a
//                            users posts (IDs) in ranked order.
//
//   [user.ID]:LIKESINORDER - KEY to ZSET containing reference keys to
//                            users liked posts IDs in chronological order.
//
//    [user.ID]:LIKESBYRANK - KEY to ZSET containing reference keys to
//                            users liked posts IDs in ranked order.
//
// [user.ID]:FRIENDSINORDER - KEY to ZSET containing reference keys to
//                            users friends IDs in chronological order.
//
//             POSTSINORDER - KEY to ZSET containing reference keys to the
//                            post IDs of every post in the database in
//                            chronological order.
//
//             POSTSBYSCORE - KEY to ZSET containing reference keys to the
//                            post IDs of every post in the database in order
//                            of rank/score.
//
// 	          [post.ID] - KEY to HASHMAP of the associated posts data.
//
// [post.ID]:REPLIESINORDER - KEY to ZSET containing reference keys to posts
//                            which are replies to other posts, in
//                            chronological order.
//
// [post.ID]:REPLIESBYSCORE - KEY to ZSET containing reference keys to posts
//                            which are replies to other posts, in ranked order
//                            (by score).
//
///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

package main

import (
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

// cache() is used in the main() function to cache the database occasionally.
// This is a crude implementation, and the final cacheing mechanism will be
// much more optimized and robust. Currently we just get the post IDs, which
// are stored in a sorted set and ranked numerically, and then get each post by
// using the keys returned, which are the post IDs in order of rank (score).
func cache() {
	postIDs, err := zrangePostsByScore() // see: zrangePostsByScore()
	if err != nil {
		log.Println(err)
		return
	}
	getPostsByID(postIDs) // see: getPostsByID()
}

// getID() is using to get the user ID associated with the password hash of a
// user whose just been authenticated. This is used to look up the users
// data/profile info.
func getID(hash string) (string, error) {
	return rdb.Get(rdx, hash).Result()
}

// zhPost(*post) is used as a one-liner to add a post to the database. This
// will add the posts ID to the ranked set, the chronological set, and add the
// post data to the database using HMSet.
func zhPost(p *post) error {
	_, err := zaddPostsChron(p) // see: zaddPostsChron()
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = zaddPostsScore(p) // see: zaddPostsScore()
	if err != nil {
		log.Println(err)
		return err
	}
	return setPost(p) // see: setPost()
}

// getPostsByID() takes a slice of IDs (use a one item slice for 1 ID) and uses
// getPost() to marshal the post data into a post{} that can be passed around
// by our program. It furthermore recursively checks each post for comments
// stored in a zset of the following pattern: post.ID:REPLIESINORDER where
// post.ID is the posts ID that which we query. Finally, we set the stream
// variable to the new []*post{}.
// TODO: Add option to sort replies by likes/score using post.ID:REPLIESBYSCORE
func getPostsByID(ids []string) []*post {
	// get the "root" level post(s).
	var items []*post = []*post{}
	for _, id := range ids {
		i, err := getPost(id) // see: getPost()
		if err != nil {
			log.Println(err)
		}
		items = append(items, &i)
	}

	// get the comments from each post. TODO: Update the amount returned
	// so it only goes a few comments deep.
	for _, p := range items {
		replies, err := rdb.ZRange(
			rdx, p.ID+":REPLIESINORDER", 0, -1,
		).Result()
		if err != nil {
			log.Println(err)
		}
		p.Comments = append(p.Comments, getPostsByID(replies)...)
	}

	// set the stream to the new slice of posts.
	stream = items

	return items
}

// setPasswordHash() is used to store the password hash in redis so that when
// a user logs in, his login name (email) can be used to look up the hash.
// Passwords are never stored in plain text.
func setPasswordHash(c *credentials, hash string) (string, error) {
	return rdb.Set(rdx, c.Name+":HASH", hash, 0).Result()
}

// getPasswordHash() is used to get the hash associated with a users login
// name, to verify the password, and to look up their ID, so that we may look
// up the user data using the hash as the key.
func getPasswordHash(c *credentials) (string, error) {
	return rdb.Get(rdx, c.Name+":HASH").Result()
}

// setHashToID() is used to look up a user ID based on the hash returned by
// getPasswordHash(). This is necessary to allow us to look up the user ID
// without needing an email, which is only used for login/verification
// purposes unless the user chooses to make it public.
func setHashToID(c *credentials, hash string) (string, error) {
	return rdb.Set(rdx, hash, c.User.ID, 0).Result()
}

// setProfile() sets a users profile data in the database by first marshalling
// it into a []byte{} containing its JSON representation, then unmarshalling it
// into a map[string]any type, which is then added to the database using the
// redis HMSet() functionality.
// TODO: There's a way to do this without a map.
func setProfile(c *credentials) error {
	// Marshal the user/profile data into its JSON representation in []byte
	// form.
	b, err := json.Marshal(c.User)
	if err != nil {
		return err
	}

	// initialize our map.
	var pmap map[string]any = make(map[string]any)

	// Unmarshal the JSON representation into the map
	err = json.Unmarshal(b, &pmap)
	if err != nil {
		return err
	}

	// Add the data using HMSet(), returning any errors.
	return rdb.HMSet(rdx, c.User.ID, pmap).Err()
}

// setPost() sets a post in the database using the redis HMSet() functionality.
func setPost(i *post) error {
	// Add the post data using HMSet(), returning any errors.
	return rdb.HMSet(rdx, i.ID, *i).Err()
}

// scanProfile() is used to scan a users profile into the &user{} struct to
// be passed around by in credentials{}.
func scanProfile(c *credentials) error {
	return rdb.HGetAll(rdx, c.User.ID).Scan(c.User)
}

// zaddUsers() is used to add a user to a sorted set called "USERS", allowing
// us to sort users by rank, which hasn't been fully implemented yet.
func zaddUsers(c *credentials) (int64, error) {
	return rdb.ZAdd(rdx, "USERS", makeZmem(c.User.ID)).Result()
}

// zaddPostsChron() is used to add a new post to the zset "POSTSINORDER", which
// maintains a chronologically sorted set of posts.
func zaddPostsChron(c *post) (int64, error) {
	return rdb.ZAdd(rdx, "POSTSINORDER", makeZmem(c.ID)).Result()
}

// zaddPostsScore() is used to add a new post to the zset "POSTSBYSCORE", which
// maintains a set of posts sorted by rank (score).
func zaddPostsScore(c *post) (int64, error) {
	return rdb.ZAdd(rdx, "POSTSBYSCORE", makeZmem(c.ID)).Result()
}

// zrangePostsByScore() returns a slice of post IDs, ordered by score, as
// returned by ZRevRangeByScore(). Obviously the naming convention is a little
// off here as tagmachine remains in testing.
// TODO: add reverse functionality.
func zrangePostsByScore() ([]string, error) {
	// TODO: Add pagification.
	opts := &redis.ZRangeBy{Min: "-inf", Max: "+inf", Offset: 0, Count: -1}
	return rdb.ZRevRangeByScore(rdx, "POSTSBYSCORE", opts).Result()
}

// getPost() is used to retrieve a single post from redis given the posts ID.
// Herein we create a &post{}, and use the redis HGetAll().Scan(&{})
// functionality, which can sometimes fill the structs key values in
// automatically, using the struct tags (and maybe best guesses).
func getPost(i string) (p post, err error) {
	return p, rdb.HGetAll(rdx, i).Scan(&p)
}

// setLike() is used to add or remove a liked post from a users liked posts
// list. It first uses ZRem to remove the post ID, but if no matching ID is
// found (returns 0), it adds the post ID, acting as a toggle-like mechanism.
// The post IDs are stored in a zset of the pattern: user.ID:LIKESINORDER
// setLike() furthermore increments or decrements a posts score accordingly, by
// looking up its post ID in the zset "POSTSBYSCORE" and using ZIncrBy() to
// increment or decrement its score, and also using HIncrByFloat() to increment
// the score stored with the post data. This may be reconfigured in the future.
// TODO: add user.ID:LIKESBYRANK sortability.
func setLike(c *credentials, id string) (int64, error) {
	// Try to remove the users like from the zset of the key pattern:
	// user.ID:LIKESINORDER, which should contain the ID's of the users
	// liked posts.
	num, err := rdb.ZRem(rdx, c.User.ID+":LIKESINORDER", -1, id).Result()
	if err != nil {
		log.Println(err)
		return -1, err
	}

	if num == 0 { // If no ID is found:

		// Add the ID to the users liked posts.
		_, err := rdb.ZAdd(rdx, c.User.ID+":LIKESINORDER",
			makeZmem(id)).Result()
		if err != nil {
			log.Println(err)
			return -1, err
		}

		// Increment the posts score in the zset of key "POSTSBYSCORE".
		_, err = rdb.ZIncrBy(rdx, "POSTSBYSCORE", 1, id).Result()
		if err != nil {
			log.Println(err)
			return -1, err
		}

		// Increment the posts score in the HSet object containing the
		// post data.
		_, err = rdb.HIncrByFloat(rdx, id, "score", 1).Result()
		if err != nil {
			log.Println(err)
			return -1, err
		}

	} else if num == 1 { // if the ID was found and removed:

		// Decrement the posts score in the zset of key "POSTSBYSCORE".
		// This is done using ZIncrBy(), with a negative value passed
		// as the increment parameter.
		_, err = rdb.ZIncrBy(rdx, "POSTSBYSCORE", -1, id).Result()
		if err != nil {
			log.Println(err)
			return -1, err
		}

		// Decrement the posts score in the HSet object containing the
		// post data. This is done using HIncrByFloat(), with a
		// negative value passed as the increment parameter.
		_, err = rdb.HIncrByFloat(rdx, id, "score", -1).Result()
		if err != nil {
			log.Println(err)
			return -1, err
		}
	} else {
		return -1, err
	}

	return num, nil
}

// setFriend() is used to add or remove a friend from a users friends list. It
// first uses ZRem to remove the friend, but if no friend is found (returns 0),
// it adds the friend, acting as a toggle-like mechanism.
// The friends IDs are stored in a set of the pattern: user.ID:FRIENDSINORDER
func setFriend(c *credentials, id string) (int64, error) {
	num, err := rdb.ZRem(rdx, c.User.ID+":FRIENDSINORDER", -1, id).Result()
	if err != nil {
		log.Println(err)
		return -1, err
	}

	if num == 0 {
		_, err := rdb.ZAdd(rdx, c.User.ID+":FRIENDSINORDER",
			makeZmem(id)).Result()
		if err != nil {
			log.Println(err)
			return -1, err
		}
	}
	return num, nil
}

// getLikes() is used to retrieve a users liked posts, stored in a zset of key
// pattern: user.ID:LIKESINORDER
func getLikes(c *credentials) ([]*post, error) {
	ids, err := rdb.ZRevRange(rdx, c.User.ID+":LIKESINORDER", 0, 10).Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return getPostsByID(ids), nil
}

// zaddUsersPosts() is used when a user submits a post. The posts ID must be
// added to the following sets in redis:
// user.ID:POSTSINORDER
// post.Parent:REPLIESINORDER
// TODO:
// user.ID:POSTSBYSCORE
// post.Parent:REPLIESBYSCORE
func zaddUsersPosts(c *credentials, p *post) (int64, error) {
	// if no parent, it's not a reply, its a root level post.
	if p.Parent == "" {
		return rdb.ZAdd(rdx, c.User.ID+":POSTSINORDER", makeZmem(p.ID)).Result()
	}

	// If a parent ID is provided via post.Parent, we look up the parent
	// post.
	pp, err := getPost(p.Parent)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	// We add the new posts ID to a sorted set containing the users post
	// IDs in chronological order.
	i, err := rdb.ZAdd(rdx, c.User.ID+":POSTSINORDER", makeZmem(p.ID)).Result()
	if err != nil {
		log.Println(err)
		return i, err
	}

	// We add the new posts ID to a sorted set containing the IDs of the
	// replies to the parent comment, so it can be looked up when the
	// parents data is queried.
	i, err = rdb.ZAdd(rdx, p.Parent+":REPLIESINORDER", makeZmem(p.ID)).Result()
	if err != nil {
		log.Println(err)
		return i, err
	}

	// Add the post data to the database.
	err = setPost(p) // see: setPost()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	// Save the updated parents data to the database.
	err = setPost(&pp) // see: setPost()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	// return 1 if all good.
	return 1, nil
}

// zaddUsersLikesChron() is used to add a users likes to a sorted set
// containing their likes inn chronological order.
func zaddUsersLikesChron(c *credentials) (int64, error) {
	return rdb.ZAdd(rdx, c.User.ID+":LIKESINORDER", makeZmem(c.User.ID)).Result()
}

// zaddTagsScore() is used to add a tag to the sorted set "TAGSBYSCORE", and
// ordered by rank.
func zaddTagsScore(c *credentials) (int64, error) {
	return rdb.ZAdd(rdx, "TAGSBYSCORE", makeZmem(c.User.ID)).Result()
}
