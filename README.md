# Go JWT Generator
Convenient tool for generating a jwt

## Build Executable
`go build main.go`

## Usage
~~~
jwt
    --alg -A        specify jwt sign alg
    --claimset -C   add claims kv
    --options -O    add extra options
    --aud -a        add audiences
    --exp           add expires at(e.g. +15s +20m +8h +24d +3M)
    --iat           add issued at(e.g. -15s +20m -8h +24d -3M)
    --sub           set subject
    --jti           set jwt ID
~~~

## Date Format
| Date Suffix | DateType                 |
|-------------|--------------------------|
| s           | time\.Second             |
| m           | time\.Minute             |
| h           | time\.Hour               |
| d           | time\.Hour x 24          |
| w           | time\.Hour x 24 x 7      |
| M           | time\.Hour x 24 x 7 x 30 |


## Claims Set and Options Parameter Format
| Value Suffix | Value Type | Example          |
|--------------|------------|------------------|
| i            | int        | 1i 23i           |
| b            | bool       | trueb falseb     |
| f            | float64    | 35\.06f 23\.4f   |
| (default)    | string     | dmao 'dmao blog' |


## Examples
~~~
Generate HS256 jwt with Claims {name:dormao,version:17,locked:true}
$ jwt -A hs256 -C name=dormao -C version=17i -C locked=trueb
~~~

## Supported Sign Algs
| Alg Name | Parameter Value |
|----------|-----------------|
| HS256    | hs256           |
| HS284    | hs384           |
| HS512    | hs512           |

## License
[MIT](https://mit-license.org/) Copyright 2020 dormao