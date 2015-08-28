// This package implements a container storing the urls which should be
// crawled in db(redis currently)
//
// Therefore, the crawler could work correct concurrently in different
// thread/process/machine.
//
// The contianer stores the url need to crawl, removes the url alreay has
// been crawled.

package urlstore
