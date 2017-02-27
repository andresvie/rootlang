# rootlang

Proof Of Concept of functional programing language 

## TODO

* Add Support to threads with gorutine
* Add Modules Support and Lambda Expression
* Add Combinators map,filter,zip,reduce
* Add Function Composition operator "."
* Add More Basic Types(float,list,map,file)
* Generate LLVM Intermediate Representation(IR)

## Syntax
rootlang has a syntax easy to follow
```rootlang
//function declaration
let x = y=>{ return y+10;};
the function has clousure by default example
let add = (x,y) => {return x + y;};
let add10 = add(10);
let x = add10(5);
//this sentences assign to variable x the value of 15, add10 became a function with the value 10 bound to local variable x in the context of the function;```
