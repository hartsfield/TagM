package main

import (
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

func zhPost(p *post) error {
	_, err := zaddPostsChron(p)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = zaddPostsScore(p)
	if err != nil {
		log.Println(err)
		return err
	}

	return setPost(p)
}
func cache() {
	postIDs, err := zrangePostsByScore()
	if err != nil {
		log.Println(err)
	}
	getPostsByID(postIDs)
}
func getPostsByID(ids []string) []*post {
	var items []*post = []*post{}
	for _, id := range ids {
		i, err := getPost(id)
		if err != nil {
			log.Println(err)
		}
		items = append(items, i)
	}
	stream = items
	return items
}
func setPasswordHash(c *credentials, hash string) (string, error) {
	return rdb.Set(rdx, c.Name+":HASH", hash, 0).Result()
}
func setHashToID(c *credentials, hash string) (string, error) {
	return rdb.Set(rdx, hash, c.User.ID, 0).Result()
}
func getPasswordHash(c *credentials) (string, error) {
	return rdb.Get(rdx, c.Name+":HASH").Result()
}
func setProfile(c *credentials) error {
	var pmap map[string]any = make(map[string]any)
	b, err := json.Marshal(c.User)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(b, &pmap)
	if err != nil {
		log.Println(err)
	}
	return rdb.HMSet(rdx, c.User.ID, pmap).Err()
	// return rdb.HSet(rdx, c.User.ID, "credentials", c).Err()
}
func scanProfile(c *credentials) error {
	return rdb.HGetAll(rdx, c.User.ID).Scan(c.User)
}
func zaddUsers(c *credentials) (int64, error) {
	return rdb.ZAdd(rdx, "USERS", makeZmem(c.User.ID)).Result()
}
func zaddPostsChron(c *post) (int64, error) {
	return rdb.ZAdd(rdx, "POSTSINORDER", makeZmem(c.ID)).Result()
}
func zaddPostsScore(c *post) (int64, error) {
	return rdb.ZAdd(rdx, "POSTSBYSCORE", makeZmem(c.ID)).Result()
}
func zaddUsersLikesChron(c *credentials) (int64, error) {
	return rdb.ZAdd(rdx, c.User.ID+"LIKESINORDER", makeZmem(c.User.ID)).Result()
}
func zaddTagsScore(c *credentials) (int64, error) {
	return rdb.ZAdd(rdx, "TAGSBYSCORE", makeZmem(c.User.ID)).Result()
}
func zrangePostsByScore() ([]string, error) {
	opts := &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  -1,
	}
	return rdb.ZRevRangeByScore(rdx, "POSTSBYSCORE", opts).Result()
}

func setPost(i *post) error {
	var pmap map[string]any = make(map[string]any)
	b, err := json.Marshal(i)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(b, &pmap)
	if err != nil {
		log.Println(err)
	}
	return rdb.HMSet(rdx, i.ID, pmap).Err()
}
func getPost(i string) (*post, error) {
	var p *post = &post{}
	err := rdb.HGetAll(rdx, i).Scan(p)
	if err != nil {
		log.Println(err)
	}
	return p, err
}
func getID(h string) (string, error) {
	return rdb.Get(rdx, h).Result()
}
func setLike(c *credentials, id string) (int64, error) {
	num, err := rdb.ZRem(rdx, c.User.ID+":LIKESINORDER", -1, id).Result()
	log.Println("----------", num)
	if err != nil {
		log.Println(err)
		return -1, err
	}

	if num == 0 {
		_, err := rdb.ZAdd(rdx, c.User.ID+":LIKESINORDER", makeZmem(id)).Result()
		if err != nil {
			log.Println(err)
			return -1, err
		}
		_, err = rdb.ZIncrBy(rdx, "POSTSBYSCORE", 1, id).Result()
		if err != nil {
			log.Println(err)
			return -1, err
		}
		_, err = rdb.HIncrByFloat(rdx, id, "score", 1).Result()
		if err != nil {
			log.Println(err)
			return -1, err
		}
	} else if num == 1 {
		_, err = rdb.ZIncrBy(rdx, "POSTSBYSCORE", -1, id).Result()
		if err != nil {
			log.Println(err)
			return -1, err
		}
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

// func getLike(c *credentials, id string) (string, error) {
// 	return rdb.get(rdx, c.User.ID+":LIKED", id, 0).Result()
// }
// func setFriend(c *credentials, hash string) (string, error) {
// 	return rdb.Set(rdx, c.Name+":HASH", hash, 0).Result()
// }
