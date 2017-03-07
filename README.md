# rootlang

Proof Of Concept of functional programing language 

## TODO

* Add Support to threads with gorutine
* Add Modules Support and Lambda Expression
* Add Function Composition operator "."
* Add More Basic Types(float,list,map,file)
* Generate LLVM Intermediate Representation(IR)

## Syntax
rootlang has a syntax easy to follow is like programming on (python,javascript and little bit of haskell) only have the good parts of them
```rootlang
//rootlang has ducktype object system a variable can has a differents values over of cycle of life;
//every valid sentences in rootlang has to be ended with semicolon character;
let x = 10;//declare integer literal bound to x variable
let x = "rootlang is awesome";// declare string literal bound to x variable
//function declaration
let x = y=>{ return y+10;};
//the function has clousure by default example
let add = (x,y) => {return x + y;};
let add10 = add(10);
let x = add10(5);
//this sentences assign to variable x the value of 15, add10 became a function with the value 10 bound to local variable x in the context of the function;
let p = list(1,2,3,4,5); // list declaration
//rootlang has support for combinators functions like map,filter,reduce,zip
let m = map(x => {return x*2;}, p); //return a new list transform by the lambda function [2,4,8,10];
let f = filter(x => {return x%2 == 0;},p); //return a new list filter by the lambda function [2,4];
let r = reduce((x,y) => {return x+y;}, p);// return a value 15 reduce by the function
let r1 = reduce((x,y) => {return x+y;}, p, 10);// return a value 25 reduce by the function with initial value of 10
let z = zip(m,f);//this return a list with another list with pair arguments [[2,2],[4,4]]
