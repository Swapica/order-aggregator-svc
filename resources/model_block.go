/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Block struct {
	Key
	Attributes BlockAttributes `json:"attributes"`
}
type BlockRequest struct {
	Data     Block    `json:"data"`
	Included Included `json:"included"`
}

type BlockListRequest struct {
	Data     []Block  `json:"data"`
	Included Included `json:"included"`
	Links    *Links   `json:"links"`
}

// MustBlock - returns Block from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustBlock(key Key) *Block {
	var block Block
	if c.tryFindEntry(key, &block) {
		return &block
	}
	return nil
}
