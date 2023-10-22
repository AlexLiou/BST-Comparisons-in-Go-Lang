/*
struct for nodes
left and right pointer to leafs
int value

Hash map for nodes


file io

arg parsing

sequential implementation
*/
package BST

import (
    "fmt"
    "flag"
    "os"
    "bufio"
    "log"
    "strings"
    "sync"
    "strconv"
    "time"
    "sort"
)

type Node struct {
    key int 
    value int 
    left *Node
    right *Node
}

type BinarySearchTree struct {
    ID int
    root *Node
}

func (bst *BinarySearchTree) Insert(key int, value int) {
    n := &Node{key, value, nil, nil}
    if bst.root == nil {
        bst.root = n
    } else {
        insertNode(bst.root, n)
    }
}

func insertNode(node, newNode *Node) {
    if newNode.value < node.value {
        if node.left == nil {
            node.left = newNode
        } else {
            insertNode(node.left, newNode)
        }
    } else {
        if node.right == nil {
            node.right = newNode
        } else {
            insertNode(node.right, newNode)
        }
    }
}

func (bst *BinarySearchTree) InOrderTraverse(f func(int)) {
    inOrderTraverse(bst.root, f)
}

func inOrderTraverse(n *Node, f func(int)) {
    if n != nil {
        inOrderTraverse(n.left, f)
        f(n.value)
        inOrderTraverse(n.right, f)
    }
}

func (bst *BinarySearchTree) String() {
    fmt.Println("--------------------")
    stringify(bst.root, 0)
    fmt.Println("--------------------")
}

func stringify(n *Node, level int) {
    if n != nil {
        format := ""
        for i := 0; i < level; i++ {
            format += "       "
        }
        format += "---["
        level++
        stringify(n.left, level)
        fmt.Printf(format+"%d\n", n.value)
        stringify(n.right, level)
    }
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func Walk(t *Node, ch chan int) {
    if t != nil {
        Walk(t.left, ch)
        ch <- t.value
        Walk(t.right, ch)
    }
}

func Walking(t *Node, ch chan int) {
    Walk(t, ch)
    defer close(ch)
}


func Same(t1, t2 *BinarySearchTree) bool {
    x := make(chan int)
    y := make(chan int)

    go Walking(t1.root, x)
    go Walking(t2.root, y)


    for {
        v1, ok1 := <-x
        v2, ok2 := <-y

        if ok1 != ok2 || v1 != v2 {
            return false
        }

        if !ok1 {
            break;
        }
    }

    return true
}

func isSameSlice(a, b []int) bool {
    if a == nil && b == nil {
        return true
    }

    if a == nil || b == nil {
        return false
    }

    if len(a) != len(b) {
        return false
    }

    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }

    return true
}

func constructTreeFromFile(inputPath string) []BinarySearchTree {
     // Read File
    file, err := os.Open("../"+inputPath)

    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    Trees := make([]BinarySearchTree, 0)

    // contruct BSTs as specified by the file
    // each line represents a different BST and ends in a newline character
    // the numbers on each line, which are separated by spaces, are the values contained in the BST
    // numbers should be inserted into the BST in the order provided
    i := 0
    for scanner.Scan() {
        numbers := strings.Fields(scanner.Text())
        bst := BinarySearchTree{i, nil}
        i++
        for index, number := range numbers {
            n, err := strconv.Atoi(number)
            if err != nil {
                log.Fatal(err)
            }
            bst.Insert(index, n)
        }
        Trees = append(Trees, bst)
    }

    return Trees
}

func computeAndPrintHashSequential(Trees []BinarySearchTree) {
    // generate a hash of each BST
    // store the hashes and identify potential duplicates

    // hash = 1;
    // for each value in tree.in_order_traversal() {
    //     new_value = value + 2;
    //     hash = (hash * new_value + new_value) % 1000
    // }
    start := time.Now()
    for _, bst := range Trees {
        hash := 1
        
        bst.InOrderTraverse(func(value int) {
            new_value := value + 2
            hash = (hash * new_value + new_value) % 1000
        })

    }

    elapsed := time.Since(start)
    // Print output of Hashmap according to guide
    // Print hashmap
    fmt.Printf("hashTime: %fs\n", elapsed.Seconds())
}

func computeAndPrintHashmapSequential(Trees []BinarySearchTree, printMap bool) map[int][]BinarySearchTree {
    // generate a hash of each BST
    // store the hashes and identify potential duplicates

    // hash = 1;
    // for each value in tree.in_order_traversal() {
    //     new_value = value + 2;
    //     hash = (hash * new_value + new_value) % 1000
    // }
    start := time.Now()
    hashmap := make(map[int][]BinarySearchTree, 0)
    hashmapPrint := make(map[int][]int, 0)
    for _, bst := range Trees {
        hash := 1
        
        bst.InOrderTraverse(func(value int) {
            new_value := value + 2
            hash = (hash * new_value + new_value) % 1000
        })

        hashmap[hash] = append(hashmap[hash], bst)
        hashmapPrint[hash] = append(hashmapPrint[hash], bst.ID)
    }

    elapsed := time.Since(start)
    // Print output of Hashmap according to guide
    // Print hashmap
    fmt.Printf("hashGroupTime: %fs\n", elapsed.Seconds())
    if printMap {
        for hash, treeIds := range hashmapPrint {
            if len(treeIds) == 1 {
                continue
            }
            fmt.Printf("hash%d ", hash)
            sort.Ints(treeIds)
            for _, id := range treeIds {
                if id < 10 {
                    fmt.Printf("id0%d ", id)
                } else {
                    fmt.Printf("id%d ", id)
                }
            }
            fmt.Println()
        }
    }

    return hashmap
}

func computeAndPrintUniqueTreesSequential(hashmap map[int][]BinarySearchTree, printMap bool) {
    // For simplicity, we recommend you use the faster of your implementations before proceeding. With your hashes parallelized, 
    // it is now time to parallelize the final tree comparisons. To keep things simple, you can store BST equivalence using a 2D 
    // adjacency matrix where array[i][j] == true implies the BST with ID i is equivalent to the BST with ID j. 
    // After you have populated your map from hashes to BST IDs and initialized the adjacency matrix to false, 
    // have a single thread traverse it. If a hash is associated with only one BST ID, k, then you can simply assign true to 
    // array[k][k] and proceed. However, if a hash has multiple BST IDs associated with it, you will have to compare all of their 
    // trees for similarity. For your first implementation, just spawn a goroutine to do each comparison and write the 
    // result to the adjacency matrix if a match is found.
    // Group Id
    groupIds := map[int][]int{}
    groupId := 0
    start := time.Now()
    for _, BSTs := range hashmap {

        // Create matrix, everything set to false
        var size = len(BSTs)
        if size == 1 {
            continue
        }
        var adjacencyMatrix = make([][]bool, size)
        for i := 0; i < size; i++ {
            adjacencyMatrix[i] = make([]bool, size)
        }

        for i := 0; i < size; i++ {
            for j := 0; j < size; j++ {
                if i < j {
                    if !adjacencyMatrix[i][j] && Same(&BSTs[i], &BSTs[j]) {
                        adjacencyMatrix[i][j] = true
                        adjacencyMatrix[j][i] = true
                    }
                }
            }
        }	


        // Figure out some way to print this into something useful
        hasBeenChecked := make([]bool, size)
        for i := 0; i < size; i++ {
            var IDs []int
            for j := 0 ; j < size; j++ {
                if i < j && adjacencyMatrix[i][j] {
                    if !hasBeenChecked[i] {
                        IDs = append(IDs, BSTs[i].ID)
                        hasBeenChecked[i] = true
                    }
                    if !hasBeenChecked[j] {
                        IDs = append(IDs, BSTs[j].ID)
                        hasBeenChecked[j] = true
                    }
                }
            }
            if len(IDs) != 0 {
                sort.Ints(IDs)
                groupIds[groupId] = IDs
                groupId++
            }
        }
    }

    elapsed := time.Since(start)
    // Print Unique Trees
    fmt.Printf("compareTreeTime: %fs\n", elapsed.Seconds())
    if printMap {
        for i := 0; i < len(groupIds); i++ {
            // Just in case
            if len(groupIds[i]) == 1 {
                continue
            }
            fmt.Printf("group %d ", i)
            for j := 0; j < len(groupIds[i]); j++ {
                treeId := groupIds[i][j]
                if treeId < 10 {
                    fmt.Printf("id0%d ", treeId)
                } else {
                    fmt.Printf("id%d ", treeId)
                }
            }
            fmt.Println()
        }
    }
}

func hashBST(bst BinarySearchTree) int {
    hash := 0
    bst.InOrderTraverse(func(value int) {
        new_value := value + 2
        hash = (hash * new_value + new_value) % 1000
    })

    return hash
}

type Result struct {
    Hash int
    Tree BinarySearchTree
}

type KeyValue struct {
    Key string
    Value int
}

// SafeSetter is safe to use concurrently.
type SafeSetter struct {
	mu sync.Mutex
	hashmap  map[int][]BinarySearchTree
    hashmapPrint map[int][]int
}

func (s *SafeSetter) safeSet(r Result) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.hashmap[r.Hash] = append(s.hashmap[r.Hash], r.Tree)
    s.hashmapPrint[r.Hash] = append(s.hashmapPrint[r.Hash], r.Tree.ID)
}

// SafeUniqueTrees is safe to use concurrently.
type SafeUniqueTrees struct {
    mu sync.Mutex
    uniqueTrees map[string][]int
}

func (s *SafeUniqueTrees) safeSet(key string, value int) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.uniqueTrees[key] = append(s.uniqueTrees[key], value)
}

func computeAndPrintHashParallel(Trees []BinarySearchTree, numHashThreads int) {
    var wg sync.WaitGroup
    
    // spawns 'hash-worer' goroutines to compute the hashes of the input BSTs. 
    start := time.Now()

    var threads int = numHashThreads
    if len(Trees) < numHashThreads {
        threads = len(Trees)
    } 
    for i := 0; i < threads; i++ {
        wg.Add(1)

        // Split Trees into 'threads' amount of subdivisions for each thread to work on.
        // i.e [0, 1, 2, 3, 4, 5, 6, 7, 8, 9] with 4 threads. 
        // Thread 1 works on 0, 1. Thread 2 works on 2, 3. Thread 3 works on 4, 5. Thread 4 works on 6, 7, 8, 9
        var step int = len(Trees) / threads
        var start int = i*step
        var end int = i*step + step
        if i == threads - 1 {
            end = len(Trees)
        }
        go func(start, end int) {
            defer wg.Done()

            for i := start; i < end; i++ {
                hashBST(Trees[i])
            }
        }(start, end)
    }

     // wait for all go routines to finish
    wg.Wait()
    elapsed := time.Since(start)
    // print hashmap
    fmt.Printf("hashTime: %fs\n", elapsed.Seconds())
}

func computeAndPrintHashmapParallel(Trees []BinarySearchTree, numHashThreads int, numDataThreads int, printMap bool) map[int][]BinarySearchTree {
    // generate a hash of each BST
    // generates a go routine for each BST
    // store the hashes and identify potential duplicates
    // Each goroutine sends its (hash, BST ID) pair(s) to a central manager goroutine using a channel.
    manager := make(chan Result)

    // create Wait Group to track go routines
    var wg sync.WaitGroup

    // numHashThreads = i, numDataThreads = 1
    if numDataThreads == 1 {
        // spawns 'hash-worer' goroutines to compute the hashes of the input BSTs. 
        start := time.Now()

        var threads int = numHashThreads
        if len(Trees) < numHashThreads {
            threads = len(Trees)
        } 
        for i := 0; i < threads; i++ {
            wg.Add(1)

            // Split Trees into 'threads' amount of subdivisions for each thread to work on.
            // i.e [0, 1, 2, 3, 4, 5, 6, 7, 8, 9] with 4 threads. 
            // Thread 1 works on 0, 1. Thread 2 works on 2, 3. Thread 3 works on 4, 5. Thread 4 works on 6, 7, 8, 9
            var step int = len(Trees) / threads
            var start int = i*step
            var end int = i*step + step
            if i == threads - 1 {
                end = len(Trees)
            }
            go func(start, end int) {
                defer wg.Done()

                for i := start; i < end; i++ {
                    hash := hashBST(Trees[i])
                    result := Result{ Hash: hash, Tree: Trees[i] }
                    manager <- result
                }
            }(start, end)
        }

        // dictionary
        hashmap := map[int][]BinarySearchTree{}
        hashmapPrint := map[int][]int{}
        // start central manager routine
        go func() {
            defer close(manager)

            // while there are still pairs in the channel, receive them and append
            for r := range manager {
                hashmap[r.Hash] = append(hashmap[r.Hash], r.Tree)
                hashmapPrint[r.Hash] = append(hashmapPrint[r.Hash], r.Tree.ID)
            }
        }()

        // wait for all go routines to finish
        wg.Wait()
        elapsed := time.Since(start)
        // print hashmap
        fmt.Printf("hashGroupTime: %fs\n", elapsed.Seconds())
        if printMap {
            for hash, treeIds := range hashmapPrint {
                if len(treeIds) == 1 {
                    continue
                }
                fmt.Printf("hash%d ", hash)
                sort.Ints(treeIds)
                for _, id := range treeIds {
                    if id < 10 {
                        fmt.Printf("id0%d ", id)
                    } else {
                        fmt.Printf("id%d ", id)
                    }
                }
                fmt.Println()
            }
        }

        return hashmap
    // Lock implementation, when hash-workers == data-workers
    } else if numHashThreads == numDataThreads {
        // spawns 'hash-worer' goroutines to compute the hashes of the input BSTs. 
        start := time.Now()

         // dictionary
        hashmap := map[int][]BinarySearchTree{}
        hashmapPrint := map[int][]int{}
        s := SafeSetter{hashmap: hashmap, hashmapPrint: hashmapPrint}

        var threads int = numHashThreads
        if len(Trees) < numHashThreads {
            threads = len(Trees)
        } 
        for i := 0; i < threads; i++ {
            wg.Add(1)

            // Split Trees into 'threads' amount of subdivisions for each thread to work on.
            // i.e [0, 1, 2, 3, 4, 5, 6, 7, 8, 9] with 4 threads. 
            // Thread 1 works on 0, 1. Thread 2 works on 2, 3. Thread 3 works on 4, 5. Thread 4 works on 6, 7, 8, 9
            var step int = len(Trees) / threads
            var start int = i*step
            var end int = i*step + step
            if i == threads - 1 {
                end = len(Trees)
            }
            go func(start, end int) {
                defer wg.Done()

                for i := start; i < end; i++ {
                    hash := hashBST(Trees[i])
                    result := Result{ Hash: hash, Tree: Trees[i] }
                    s.safeSet(result)
                }
            }(start, end)
        }

        // wait for all go routines to finish
        wg.Wait()
        elapsed := time.Since(start)
        // print hashmap
        fmt.Printf("hashGroupTime: %fs\n", elapsed.Seconds())
        if printMap {
            for hash, treeIds := range s.hashmapPrint {
                if len(treeIds) == 1 {
                    continue
                }
                fmt.Printf("hash%d ", hash)
                sort.Ints(treeIds)
                for _, id := range treeIds {
                    if id < 10 {
                        fmt.Printf("id0%d ", id)
                    } else {
                        fmt.Printf("id%d ", id)
                    }
                }
                fmt.Println()
            }
        }

        return s.hashmap

    } else {
        fmt.Println("This is for the optional implementation and thus unimplemented. Sorry.")
        return map[int][]BinarySearchTree{}
    }
}

type BSTPair struct {
    a BinarySearchTree
    b BinarySearchTree
    index_a int
    index_b int
}

type ConcurrentBuffer struct {
	data     []BSTPair
	maxSize  int
	mu       sync.Mutex
	notFull  *sync.Cond
	notEmpty *sync.Cond
}

func NewConcurrentBuffer(maxSize int) *ConcurrentBuffer {
	buffer := &ConcurrentBuffer{
		data:    make([]BSTPair, 0, maxSize),
		maxSize: maxSize,
	}
	buffer.notFull = sync.NewCond(&buffer.mu)
	buffer.notEmpty = sync.NewCond(&buffer.mu)
	return buffer
}

func (b *ConcurrentBuffer) Insert(data BSTPair) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for len(b.data) == b.maxSize {
		b.notFull.Wait()
	}
	b.data = append(b.data, data)
	b.notEmpty.Signal()
}

func (b *ConcurrentBuffer) Remove() BSTPair {
	b.mu.Lock()
	defer b.mu.Unlock()
	for len(b.data) == 0 {
		b.notEmpty.Wait()
	}
	data := b.data[0]
	b.data = b.data[1:]
	b.notFull.Signal()
	return data
}

func computeAndPrintUniqueTreesAdjacency(hashmap map[int][]BinarySearchTree, numCompThreads int, printMap bool) {

    // For your second implementation, you will spawn comp-workers threads to do the comparisons and use a concurrent buffer to 
    // communicate with them (work can be represented with a (BST ID, BST ID) pair of trees to compare). 
    // Since this is not a data structures course, it doesn't matter how the buffer is implemented as long as you make sure 
    // that it only holds up to comp-workers items at a time and isn't slow. Your buffer should contain a mutex to prevent 
    // concurrency errors when multiple threads try to access the buffer at once. The buffer should also contain two 
    // conditions to handle the cases where the main thread tries to insert work when the buffer is full and when the 
    // worker threads try to remove work when the buffer is empty. The threads should not "spin" and repeatedly check 
    // if the buffer is no longer empty or full.
    
    // Create a concurrent buffer with n worker threads
    var wg sync.WaitGroup
    buffer := NewConcurrentBuffer(numCompThreads)
    // Group Id
    groupIds := map[int][]int{}
    groupId := 0
    start := time.Now()
    for _, BSTs := range hashmap {

        // Create matrix, everything set to false
        var size = len(BSTs)
        if size == 1 {
            continue
        }
        var adjacencyMatrix = make([][]bool, size)
        for i := 0; i < size; i++ {
            adjacencyMatrix[i] = make([]bool, size)
        }
        // if Same(&pair.a, &pair.b) {
        //     adjacencyMatrix[pair.index_a][pair.index_b] = true
        // }

        // Spawn worker threads.
        exit := make(chan bool)

        threads := numCompThreads
        if size < numCompThreads {
            threads = size
        }

        // Because we know exactly how many insertions and removes need to be done, we can 'cheat' and set hard limits
        // on the amount of times the go routines run. Additionally, because each operation accesses a different
        // of the array, there's no need to worry about race conditions for the adjacencyMatrix
        // Create n remove goroutines
        for i := 0; i < threads; i++ {
            wg.Add(1)
            go func(id int) {
                var step int = size*size / threads
                var start int = i*step
                var end int = i*step + step
                if i == threads - 1 {
                    end = size*size
                }
                defer wg.Done()
                for i := start; i < end; i++ {
                    pair := buffer.Remove()
                    if !adjacencyMatrix[pair.index_a][pair.index_b] && Same(&pair.a, &pair.b) {
                        adjacencyMatrix[pair.index_a][pair.index_b] = true
                        adjacencyMatrix[pair.index_b][pair.index_a] = true
                    }
                    // fmt.Printf("Removed by goroutine %d: (%d, %d)\n", id, pair.a.ID, pair.b.ID)
                }
            }(i)
        }

        go func() {
            wg.Add(1)
            defer wg.Done()
            for i := 0; i < size; i++ {
                for j := 0; j < size; j++ {
                    data := BSTPair{a: BSTs[i], b: BSTs[j], index_a: i, index_b: j}
                    buffer.Insert(data)
                }
            }	
            close(exit)
        }()
        
        wg.Wait()

        // Figure out some way to print this into something useful
        hasBeenChecked := make([]bool, size)
        for i := 0; i < size; i++ {
            var IDs []int
            for j := 0 ; j < size; j++ {
                if i < j && adjacencyMatrix[i][j] {
                    if !hasBeenChecked[i] {
                        IDs = append(IDs, BSTs[i].ID)
                        hasBeenChecked[i] = true
                    }
                    if !hasBeenChecked[j] {
                        IDs = append(IDs, BSTs[j].ID)
                        hasBeenChecked[j] = true
                    }
                }
            }
            if len(IDs) != 0 {
                sort.Ints(IDs)
                groupIds[groupId] = IDs
                groupId++
            }
        }
    }

    elapsed := time.Since(start)
    // Print Unique Trees
    fmt.Printf("compareTreeTime: %fs\n", elapsed.Seconds())
    if printMap {
        for i := 0; i < len(groupIds); i++ {
            // Just in case
            if len(groupIds[i]) == 1 {
                continue
            }
            fmt.Printf("group %d ", i)
            for j := 0; j < len(groupIds[i]); j++ {
                treeId := groupIds[i][j]
                if treeId < 10 {
                    fmt.Printf("id0%d ", treeId)
                } else {
                    fmt.Printf("id%d ", treeId)
                }
            }
            fmt.Println()
        }
    }
}

func computeAndPrintUniqueTreesParallelLock(hashmap map[int][]BinarySearchTree, numCompThreads int, printMap bool) {
    // compare trees with identical hashes to determine their equality
    // for each tree, store the IDs of each identical tree
    // a BST's ID is just the index of it in the input file

    // Take list of trees. Check if list is of size 1. Ignore if true.
    // Take first tree. Check if each subsequent tree is same as first tree.
    // If true, remove it from the list and add it to the uniqueTrees. If false, ignore.
    // create Wait Group to track go routines
     // Use Lock for simplicity and performance

    uniqueTrees := map[string][]int{}

    var wg sync.WaitGroup
    s := SafeUniqueTrees{uniqueTrees: uniqueTrees}
    start := time.Now()
    for _, trees := range hashmap {
        if len(trees) == 1 {
            continue
        }

        var threads int = numCompThreads
        if len(trees) < numCompThreads {
            threads = len(trees)
        } 
        for i := 0; i < threads; i++ {
            wg.Add(1)

            // Split Trees into 'threads' amount of subdivisions for each thread to work on.
            // i.e [0, 1, 2, 3, 4, 5, 6, 7, 8, 9] with 4 threads. 
            // Thread 1 works on 0, 1. Thread 2 works on 2, 3. Thread 3 works on 4, 5. Thread 4 works on 6, 7, 8, 9
            // [0, 1], step = 1, start = 0 and 1, end = 1 and 2
            // [0], step = 1, start 
            var step int = len(trees) / threads
            var start int = i*step
            var end int = i*step + step

            if i == threads - 1 {
                end = len(trees)
            }
            
            go func(start, end int, t []BinarySearchTree) {
                defer wg.Done()

                for i := start; i < end; i++ {
                    var result []string
                    // fmt.Println(len(t), i)
                    t[i].InOrderTraverse(func(i int) {
                        result = append(result, strconv.Itoa(i))
                    })
                    key := strings.Join(result, "#")
                    s.safeSet(key, t[i].ID)
                }
            }(start, end, trees)
        }

        
    }

    groupId := 0
    // wait for all go routines to finish
    wg.Wait()
    elapsed := time.Since(start)
    // Print Unique Trees
    fmt.Printf("compareTreeTime: %fs\n", elapsed.Seconds())
    if printMap {
        for _, treeIds := range uniqueTrees {
            if len(treeIds) == 1 {
                continue
            }
            fmt.Printf("group %d ", groupId)
            sort.Ints(treeIds)
            for _, treeId := range treeIds {
                if treeId < 10 {
                    fmt.Printf("id0%d ", treeId)
                } else {
                    fmt.Printf("id%d ", treeId)
                }
            }
            fmt.Println()
            groupId++
        }   
    }
}

func main() {
    var hashThreads int
    var dataThreads int
    var compThreads int
    var inputPath string
    var printMap = true
    // Get and Set Command-line arguments
    // -hash-workers= int number of threads
    flag.IntVar(&hashThreads, "hash-workers", 0, "number of threads for hash")
    // -data-workers= int number of threads
    flag.IntVar(&dataThreads, "data-workers", 0, "number of threads for data")
    // -comp-workers= int number of threads
    flag.IntVar(&compThreads, "comp-workers", 0, "number of threads for comp")
    // -input= string path to an input file
    flag.StringVar(&inputPath, "input", "", "path to an input file")
    
    flag.Parse()


    if hashThreads < 0 || dataThreads < 0 || compThreads < 0 {
        fmt.Println("Invalid amount of threads. Threads can't be negative.")
        return
    }
    
    Trees := constructTreeFromFile(inputPath)

    if dataThreads == 0 && compThreads == 0 {

        if hashThreads == 1 {
            computeAndPrintHashSequential(Trees)
        } else {
            computeAndPrintHashParallel(Trees, hashThreads)
        }
        return
    }

    if hashThreads != 0 && dataThreads != 0 {
        hashmap := map[int][]BinarySearchTree{}

        if hashThreads == 1 {
            hashmap = computeAndPrintHashmapSequential(Trees, printMap)
        } else {
            hashmap = computeAndPrintHashmapParallel(Trees, hashThreads, dataThreads, printMap)
        }

        if compThreads != 0 {

            if compThreads == 1 {
                computeAndPrintUniqueTreesSequential(hashmap, printMap)
            } else {
                computeAndPrintUniqueTreesAdjacency(hashmap, compThreads, printMap)
            }
        }
    }
    return 
}