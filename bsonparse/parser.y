%{
package bsonparse

type pair struct {
  key string
  val interface{}
}

func setResult(l yyLexer, v []interface{}) {
  l.(*lex).result = v
}
%}

%union{
  obj map[string]interface{}
  list []interface{}
  pair pair
  val interface{}
}

%token LexError
%token BsonError
%token <val> String Number Literal ObjectID ISODate NumberLong NumberDecimal Undefined MinKey MaxKey DBRef BinData

%type <obj> object members
%type <pair> pair
%type <val> array
%type <list> elements
%type <val> value


%start array

%%
array: '[' elements ']'
  {
    $$ = $2
    setResult(yylex, $2)
  }

object: '{' members '}'
  {
    $$ = $2
  }

members:
  {
    $$ = map[string]interface{}{}
  }
| pair
  {
    $$ = map[string]interface{}{
      $1.key: $1.val,
    }
  }
| members ',' pair
  {
    $1[$3.key] = $3.val
    $$ = $1
  }

pair:
String ':' value
  {
    $$ = pair{key: $1.(string), val: $3}
  }
| '"' String '"' ':' value
  {
    $$ = pair{key: $2.(string), val: $5}
  }

elements:
  {
    $$ = []interface{}{}
  }
| value
  {
    $$ = []interface{}{$1}
  }
| elements ',' value
  {
    $$ = append($1, $3)
  }

value:
Number
| '"' '"'
  {
    $$ = ""
  }
| '"' String '"'
  {
    $$ = $2
  }
| Literal
| object
  {
    $$ = $1
  }
| array
| ObjectID '(' '"' String '"' ')'
  {
    $$ = map[string]interface{}{"$oid": $4}
  }
| ISODate '(' '"' String '"' ')'
  {
    $$ = map[string]interface{}{"$date": $4}
  }
| NumberLong '(' '"' String '"' ')'
  {
    $$ = map[string]interface{}{"$numberLong": $4}
  }
| NumberDecimal '(' '"' String '"' ')'
  {
    $$ = map[string]interface{}{"$numberDecimal": $4}
  }
| Undefined
  {
    $$ = map[string]interface{}{"$undefined": true}
  }
| MinKey
  {
    $$ = map[string]interface{}{"$minKey": true}
  }
| MaxKey
  {
    $$ = map[string]interface{}{"$maxKey": true}
  }
| DBRef '(' '"' String '"' ',' '"' String '"' ')'
  {
    $$ = map[string]interface{}{"$ref": $4, "$id": $8}
  }
| BinData '(' String ',' String ')'
  {
    $$ = map[string]interface{}{"$binary": $5, "$type": $3}
  }
| '/' String '/' String
  {
    $$ = map[string]interface{}{"$regex": $2, "$options": $4}
  }