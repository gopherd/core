# API Reference

<ul>
  <li><a href="#user-content-_">_</a></li>
  <li><a href="#user-content-Convert">Convert</a></li>
<ul>
    <li><a href="#user-content-Convert_bool">bool</a></li>
    <li><a href="#user-content-Convert_float">float</a></li>
    <li><a href="#user-content-Convert_int">int</a></li>
    <li><a href="#user-content-Convert_string">string</a></li>
</ul>
  <li><a href="#user-content-Date">Date</a></li>
<ul>
    <li><a href="#user-content-Date_now">now</a></li>
    <li><a href="#user-content-Date_parseTime">parseTime</a></li>
</ul>
  <li><a href="#user-content-Encoding">Encoding</a></li>
<ul>
    <li><a href="#user-content-Encoding_b64dec">b64dec</a></li>
    <li><a href="#user-content-Encoding_b64enc">b64enc</a></li>
</ul>
  <li><a href="#user-content-List">List</a></li>
<ul>
    <li><a href="#user-content-List_first">first</a></li>
    <li><a href="#user-content-List_includes">includes</a></li>
    <li><a href="#user-content-List_last">last</a></li>
    <li><a href="#user-content-List_list">list</a></li>
    <li><a href="#user-content-List_map">map</a></li>
    <li><a href="#user-content-List_reverse">reverse</a></li>
    <li><a href="#user-content-List_sort">sort</a></li>
    <li><a href="#user-content-List_uniq">uniq</a></li>
</ul>
  <li><a href="#user-content-Math">Math</a></li>
<ul>
    <li><a href="#user-content-Math_add">add</a></li>
    <li><a href="#user-content-Math_ceil">ceil</a></li>
    <li><a href="#user-content-Math_floor">floor</a></li>
    <li><a href="#user-content-Math_max">max</a></li>
    <li><a href="#user-content-Math_min">min</a></li>
    <li><a href="#user-content-Math_mod">mod</a></li>
    <li><a href="#user-content-Math_mul">mul</a></li>
    <li><a href="#user-content-Math_quo">quo</a></li>
    <li><a href="#user-content-Math_rem">rem</a></li>
    <li><a href="#user-content-Math_round">round</a></li>
    <li><a href="#user-content-Math_sub">sub</a></li>
</ul>
  <li><a href="#user-content-Strings">Strings</a></li>
<ul>
    <li><a href="#user-content-Strings_camelCase">camelCase</a></li>
    <li><a href="#user-content-Strings_capitalize">capitalize</a></li>
    <li><a href="#user-content-Strings_center">center</a></li>
    <li><a href="#user-content-Strings_hasPrefix">hasPrefix</a></li>
    <li><a href="#user-content-Strings_hasSuffix">hasSuffix</a></li>
    <li><a href="#user-content-Strings_html">html</a></li>
    <li><a href="#user-content-Strings_join">join</a></li>
    <li><a href="#user-content-Strings_kebabCase">kebabCase</a></li>
    <li><a href="#user-content-Strings_lower">lower</a></li>
    <li><a href="#user-content-Strings_matchRegex">matchRegex</a></li>
    <li><a href="#user-content-Strings_pascalCase">pascalCase</a></li>
    <li><a href="#user-content-Strings_quote">quote</a></li>
    <li><a href="#user-content-Strings_repeat">repeat</a></li>
    <li><a href="#user-content-Strings_replace">replace</a></li>
    <li><a href="#user-content-Strings_replaceN">replaceN</a></li>
    <li><a href="#user-content-Strings_snakeCase">snakeCase</a></li>
    <li><a href="#user-content-Strings_split">split</a></li>
    <li><a href="#user-content-Strings_striptags">striptags</a></li>
    <li><a href="#user-content-Strings_substr">substr</a></li>
    <li><a href="#user-content-Strings_trim">trim</a></li>
    <li><a href="#user-content-Strings_trimPrefix">trimPrefix</a></li>
    <li><a href="#user-content-Strings_trimSuffix">trimSuffix</a></li>
    <li><a href="#user-content-Strings_truncate">truncate</a></li>
    <li><a href="#user-content-Strings_unquote">unquote</a></li>
    <li><a href="#user-content-Strings_upper">upper</a></li>
    <li><a href="#user-content-Strings_urlEscape">urlEscape</a></li>
    <li><a href="#user-content-Strings_urlUnescape">urlUnescape</a></li>
    <li><a href="#user-content-Strings_wordwrap">wordwrap</a></li>
</ul>
</ul>

<h2><a id="user-content-_" target="_self">_</a></h2>

`_` is a no-op function that returns an empty string. It's useful to place a newline in the template. 
Example: 

```
{{- if .Ok}}
{{printf "ok: %v" .Ok}}
{{_}}
{{- end}}
```

<h2><a id="user-content-Convert" target="_self">Convert</a></h2>
<h3><a id="user-content-Convert_bool" target="_self">bool</a></h3>

`bool` converts a value to a boolean. 
Example: 
```
{{bool 1}}
{{bool "false"}}
```

Output: 
```
true
false
```

<h3><a id="user-content-Convert_float" target="_self">float</a></h3>

`float` converts a value to a float. 
Example: 
```
{{float "3.14"}}
{{float 42}}
```

Output: 
```
3.14
42
```

<h3><a id="user-content-Convert_int" target="_self">int</a></h3>

`int` converts a value to an integer. 
Example: 
```
{{int "42"}}
{{int 3.14}}
```

Output: 
```
42
3
```

<h3><a id="user-content-Convert_string" target="_self">string</a></h3>

`string` converts a value to a string. 
Example: 
```
{{string 42}}
{{string true}}
```

Output: 
```
42
true
```

<h2><a id="user-content-Date" target="_self">Date</a></h2>
<h3><a id="user-content-Date_now" target="_self">now</a></h3>

`now` returns the current time. 
Example: 
```
{{now}}
```

Output: 
```
2024-09-12 15:04:05.999999999 +0000 UTC
```

<h3><a id="user-content-Date_parseTime" target="_self">parseTime</a></h3>

`parseTime` parses a time string using the specified layout. 
- **Parameters**: (_layout_: string, _value_: string)

Example: 
```
{{parseTime "2006-01-02" "2024-09-12"}}
```

Output: 
```
2024-09-12 00:00:00 +0000 UTC
```

<h2><a id="user-content-Encoding" target="_self">Encoding</a></h2>
<h3><a id="user-content-Encoding_b64dec" target="_self">b64dec</a></h3>

`b64dec` decodes a base64 encoded string. 
Example: 
```
{{b64dec "SGVsbG8sIFdvcmxkIQ=="}}
```

Output: 
```
Hello, World!
```

<h3><a id="user-content-Encoding_b64enc" target="_self">b64enc</a></h3>

`b64enc` encodes a string to base64. 
Example: 
```
{{b64enc "Hello, World!"}}
```

Output: 
```
SGVsbG8sIFdvcmxkIQ==
```

<h2><a id="user-content-List" target="_self">List</a></h2>
<h3><a id="user-content-List_first" target="_self">first</a></h3>

`first` returns the first element of a list or string. 
Example: 
```
{{first (list 1 2 3)}}
{{first "hello"}}
```

Output: 
```
1
h
```

<h3><a id="user-content-List_includes" target="_self">includes</a></h3>

`includes` checks if an item is present in a list, map, or string. 
- **Parameters**: (_item_: any, _collection_: slice | map | string)
- **Returns**: bool

Example: 
```
{{includes 2 (list 1 2 3)}}
{{includes "world" "hello world"}}
```

Output: 
```
true
true
```

<h3><a id="user-content-List_last" target="_self">last</a></h3>

`last` returns the last element of a list or string. 
Example: 
```
{{last (list 1 2 3)}}
{{last "hello"}}
```

Output: 
```
3
o
```

<h3><a id="user-content-List_list" target="_self">list</a></h3>

`list` creates a list from the given arguments. 
Example: 
```
{{list 1 2 3}}
```

Output: 
```
[1 2 3]
```

<h3><a id="user-content-List_map" target="_self">map</a></h3>

`map` maps a list of values using the given function and returns a list of results. 
- **Parameters**: (_fn_: function, _list_: slice)

Example: 
```
{{list 1 2 3 | map (add 1)}}
{{list "a" "b" "c" | map (upper | replace "A" "X")}}
```

Output: 
```
[2 3 4]
[X B C]
```

<h3><a id="user-content-List_reverse" target="_self">reverse</a></h3>

`reverse` reverses a list or string. 
Example: 
```
{{reverse (list 1 2 3)}}
{{reverse "hello"}}
```

Output: 
```
[3 2 1]
olleh
```

<h3><a id="user-content-List_sort" target="_self">sort</a></h3>

`sort` sorts a list of numbers or strings. 
Example: 
```
{{sort (list 3 1 4 1 5 9)}}
{{sort (list "banana" "apple" "cherry")}}
```

Output: 
```
[1 1 3 4 5 9]
[apple banana cherry]
```

<h3><a id="user-content-List_uniq" target="_self">uniq</a></h3>

`uniq` removes duplicate elements from a list. 
Example: 
```
{{uniq (list 1 2 2 3 3 3)}}
```

Output: 
```
[1 2 3]
```

<h2><a id="user-content-Math" target="_self">Math</a></h2>
<h3><a id="user-content-Math_add" target="_self">add</a></h3>

`add` adds two numbers. 
- **Parameters**: (_a_: number, _b_: number)

Example: 
```
{{add 2 3}}
```

Output: 
```
5
```

<h3><a id="user-content-Math_ceil" target="_self">ceil</a></h3>

`ceil` returns the least integer value greater than or equal to the input. 
Example: 
```
{{ceil 3.14}}
```

Output: 
```
4
```

<h3><a id="user-content-Math_floor" target="_self">floor</a></h3>

`floor` returns the greatest integer value less than or equal to the input. 
Example: 
```
{{floor 3.14}}
```

Output: 
```
3
```

<h3><a id="user-content-Math_max" target="_self">max</a></h3>

`max` returns the maximum of a list of numbers. 
- **Parameters**: numbers (variadic)

Example: 
```
{{max 3 1 4 1 5 9}}
```

Output: 
```
9
```

<h3><a id="user-content-Math_min" target="_self">min</a></h3>

`min` returns the minimum of a list of numbers. 
- **Parameters**: numbers (variadic)

Example: 
```
{{min 3 1 4 1 5 9}}
```

Output: 
```
1
```

<h3><a id="user-content-Math_mod" target="_self">mod</a></h3>

`mod` returns the modulus of dividing the first number by the second. 
- **Parameters**: (_a_: number, _b_: number)

Example: 
```
{{mod -7 3}}
```

Output: 
```
2
```

<h3><a id="user-content-Math_mul" target="_self">mul</a></h3>

`mul` multiplies two numbers. 
- **Parameters**: (_a_: number, _b_: number)

Example: 
```
{{mul 2 3}}
```

Output: 
```
6
```

<h3><a id="user-content-Math_quo" target="_self">quo</a></h3>

`quo` divides the first number by the second. 
- **Parameters**: (_a_: number, _b_: number)

Example: 
```
{{quo 6 3}}
```

Output: 
```
2
```

<h3><a id="user-content-Math_rem" target="_self">rem</a></h3>

`rem` returns the remainder of dividing the first number by the second. 
- **Parameters**: (_a_: number, _b_: number)

Example: 
```
{{rem 7 3}}
```

Output: 
```
1
```

<h3><a id="user-content-Math_round" target="_self">round</a></h3>

`round` rounds a number to a specified number of decimal places. 
- **Parameters**: (_precision_: integer, _value_: number)

Example: 
```
{{round 2 3.14159}}
```

Output: 
```
3.14
```

<h3><a id="user-content-Math_sub" target="_self">sub</a></h3>

`sub` subtracts the second number from the first. 
- **Parameters**: (_a_: number, _b_: number)

Example: 
```
{{sub 5 3}}
```

Output: 
```
2
```

<h2><a id="user-content-Strings" target="_self">Strings</a></h2>
<h3><a id="user-content-Strings_camelCase" target="_self">camelCase</a></h3>

`camelCase` converts a string to camelCase. 
Example: 
```
{{camelCase "hello world"}}
```

Output: 
```
helloWorld
```

<h3><a id="user-content-Strings_capitalize" target="_self">capitalize</a></h3>

`capitalize` capitalizes the first character of a string. 
Example: 
```
{{capitalize "hello"}}
```

Output: 
```
Hello
```

<h3><a id="user-content-Strings_center" target="_self">center</a></h3>

`center` centers a string in a field of a given width. 
- **Parameters**: (_width_: int, _target_: string)

Example: 
```
{{center 20 "Hello"}}
```

Output: 
```
"       Hello        "
```

<h3><a id="user-content-Strings_hasPrefix" target="_self">hasPrefix</a></h3>

`hasPrefix` checks if a string starts with a given prefix. 
- **Parameters**: (_prefix_: string, _target_: string)
- **Returns**: bool

Example: 
```
{{hasPrefix "Hello" "Hello, World!"}}
```

Output: 
```
true
```

<h3><a id="user-content-Strings_hasSuffix" target="_self">hasSuffix</a></h3>

`hasSuffix` checks if a string ends with a given suffix. 
- **Parameters**: (_suffix_: string, _target_: string)
- **Returns**: bool

Example: 
```
{{hasSuffix "World!" "Hello, World!"}}
```

Output: 
```
true
```

<h3><a id="user-content-Strings_html" target="_self">html</a></h3>

`html` escapes special characters in a string for use in HTML. 
Example: 
```
{{html "<script>alert('XSS')</script>"}}
```

Output: 
```
&lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;
```

<h3><a id="user-content-Strings_join" target="_self">join</a></h3>

`join` joins a slice of strings with a separator. 
- **Parameters**: (_separator_: string, _values_: slice of strings)
- **Returns**: string

Example: 
```
{{join "-" (list "apple" "banana" "cherry")}}
```

Output: 
```
apple-banana-cherry
```

<h3><a id="user-content-Strings_kebabCase" target="_self">kebabCase</a></h3>

`kebabCase` converts a string to kebab-case. 
Example: 
```
{{kebabCase "helloWorld"}}
```

Output: 
```
hello-world
```

<h3><a id="user-content-Strings_lower" target="_self">lower</a></h3>

`lower` converts a string to lowercase. 
Example: 
```
{{lower "HELLO"}}
```

Output: 
```
hello
```

<h3><a id="user-content-Strings_matchRegex" target="_self">matchRegex</a></h3>

`matchRegex` checks if a string matches a regular expression. 
- **Parameters**: (_pattern_: string, _target_: string)
- **Returns**: bool

Example: 
```
{{matchRegex "^[a-z]+$" "hello"}}
```

Output: 
```
true
```

<h3><a id="user-content-Strings_pascalCase" target="_self">pascalCase</a></h3>

`pascalCase` converts a string to PascalCase. 
Example: 
```
{{pascalCase "hello world"}}
```

Output: 
```
HelloWorld
```

<h3><a id="user-content-Strings_quote" target="_self">quote</a></h3>

`quote` returns a double-quoted string. 
Example: 
```
{{print "hello"}}
{{quote "hello"}}
```

Output: 
```
hello
"hello"
```

<h3><a id="user-content-Strings_repeat" target="_self">repeat</a></h3>

`repeat` repeats a string a specified number of times. 
- **Parameters**: (_count_: int, _target_: string)

Example: 
```
{{repeat 3 "abc"}}
```

Output: 
```
abcabcabc
```

<h3><a id="user-content-Strings_replace" target="_self">replace</a></h3>

`replace` replaces all occurrences of a substring with another substring. 
- **Parameters**: (_old_: string, _new_: string, _target_: string)

Example: 
```
{{replace "o" "0" "hello world"}}
```

Output: 
```
hell0 w0rld
```

<h3><a id="user-content-Strings_replaceN" target="_self">replaceN</a></h3>

`replaceN` replaces the first n occurrences of a substring with another substring. 
- **Parameters**: (_old_: string, _new_: string, _n_: int, _target_: string)

Example: 
```
{{replaceN "o" "0" 1 "hello world"}}
```

Output: 
```
hell0 world
```

<h3><a id="user-content-Strings_snakeCase" target="_self">snakeCase</a></h3>

`snakeCase` converts a string to snake_case. 
Example: 
```
{{snakeCase "helloWorld"}}
```

Output: 
```
hello_world
```

<h3><a id="user-content-Strings_split" target="_self">split</a></h3>

`split` splits a string by a separator. 
- **Parameters**: (_separator_: string, _target_: string)
- **Returns**: slice of strings

Example: 
```
{{split "," "apple,banana,cherry"}}
```

Output: 
```
[apple banana cherry]
```

<h3><a id="user-content-Strings_striptags" target="_self">striptags</a></h3>

`striptags` removes HTML tags from a string. 
Example: 
```
{{striptags "<p>Hello <b>World</b>!</p>"}}
```

Output: 
```
Hello World!
```

<h3><a id="user-content-Strings_substr" target="_self">substr</a></h3>

`substr` extracts a substring from a string. 
- **Parameters**: (_start_: int, _length_: int, _target_: string)

Example: 
```
{{substr 0 5 "Hello, World!"}}
```

Output: 
```
Hello
```

<h3><a id="user-content-Strings_trim" target="_self">trim</a></h3>

`trim` removes leading and trailing whitespace from a string. 
Example: 
```
{{trim "  hello  "}}
```

Output: 
```
hello
```

<h3><a id="user-content-Strings_trimPrefix" target="_self">trimPrefix</a></h3>

`trimPrefix` removes a prefix from a string if it exists. 
- **Parameters**: (_prefix_: string, _target_: string)

Example: 
```
{{trimPrefix "Hello, " "Hello, World!"}}
```

Output: 
```
World!
```

<h3><a id="user-content-Strings_trimSuffix" target="_self">trimSuffix</a></h3>

`trimSuffix` removes a suffix from a string if it exists. 
- **Parameters**: (_suffix_: string, _target_: string)

Example: 
```
{{trimSuffix ", World!" "Hello, World!"}}
```

Output: 
```
Hello
```

<h3><a id="user-content-Strings_truncate" target="_self">truncate</a></h3>

`truncate` truncates a string to a specified length and adds a suffix if truncated. 
- **Parameters**: (_length_: int, _suffix_: string, _target_: string)

Example: 
```
{{truncate 10 "..." "This is a long sentence."}}
```

Output: 
```
This is a...
```

<h3><a id="user-content-Strings_unquote" target="_self">unquote</a></h3>

`unquote` returns an unquoted string. 
Example: 
```
{{unquote "\"hello\""}}
```

Output: 
```
hello
```

<h3><a id="user-content-Strings_upper" target="_self">upper</a></h3>

`upper` converts a string to uppercase. 
Example: 
```
{{upper "hello"}}
```

Output: 
```
HELLO
```

<h3><a id="user-content-Strings_urlEscape" target="_self">urlEscape</a></h3>

`urlEscape` escapes a string for use in a URL query. 
Example: 
```
{{urlEscape "hello world"}}
```

Output: 
```
hello+world
```

<h3><a id="user-content-Strings_urlUnescape" target="_self">urlUnescape</a></h3>

`urlUnescape` unescapes a URL query string. 
Example: 
```
{{urlUnescape "hello+world"}}
```

Output: 
```
hello world
```

<h3><a id="user-content-Strings_wordwrap" target="_self">wordwrap</a></h3>

`wordwrap` wraps words in a string to a specified width. 
- **Parameters**: (_width_: int, _target_: string)

Example: 
```
{{wordwrap 10 "This is a long sentence that needs wrapping."}}
```

Output: 
```
This is a
long
sentence
that needs
wrapping.
```

