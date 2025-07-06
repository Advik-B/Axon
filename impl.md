# Some implementation rules for the Axon Graph

## Execution Flow

- Any start node must eventually connect to an end node (execution flow)
- Global constants and global variables, and global types (structs) get special treatment, they don't have to be included in the Execution Flow but when transpiled, they will be declared after the imports section (special treatment for globals)
- The longest path/chain of execution flow will be the main function (will end in } ), if END is called in a specific branch, it should be transpiled into a return statement
- Enforce every path of the execution flow that started from a start/start-like(future feature) node will end in the END node.


## Special treatment of globals

- Also applies to functions, basically global-anything are excepmpt from execution flow, BUT function CALLING must adhere to [Execution Flow]
- Any and all imports the programs makes should be defined in the graph file itself

## Start-Likes
- Essentially, there are used to define functions in the program. these functions can be local or global.
- Start-like can be used to accept/require parameters for functions to work.
- There is no end-like node. but you can end with value(s), this will translate into the return statement

## Structs
- Structs can be defined independent of execution flow or within the path of execution
- Structs that are defined independent of execution flow will be put after the imports
- Axon goes out of its way to forbid in-line or "anonymus" structs. Although this will force the programmer to define all structs explictly. This rule will make debugging much easier
```go

return struct{
    name string
    age int
}{"Gopher", 17} // Transpilations like this should not be allowed

```

## Methods for structs
- Same as functions but associated with a struct
- Need to be able to differenciate when a struct is used as a refernce vs when a sturct is copied

```go

type S {
    id string
    value bool
}

// The user should be able to differenciate between
// this
func (s S) String() string {
    return s.id + string(s.Value)
}

// And this
func (s *S) setValue(val bool) {
    s.value = val
}
```

## Comments

Comments can be plaintext or markdown.
- For a comment to be included in the transpilation it must be attached to a node 
- A comment may be attached to multiple nodes (when transpiled, the comment will be copied over (be)for every node)

### Comments (special case)
- Floating/Independent  comments
- Comments can also be not associated with any node (although this is not recomended, we dont disallow it)
- Comments like these are considered purely visual (will only show up on the node graph. but not the editor)
- There comments will be LOST on transpilation to go code


## Visual Info/Attributes
- Nodes can optionally have a visual attribute, when defines their xy position in graph and also their dimentions
- This is purely cosmetic and does not affect the go code