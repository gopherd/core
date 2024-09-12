# API Reference

<ul>
  <li><a href="#user-content-Function">Function</a></li>
<ul>
    <li><a href="#user-content-Function__">_</a></li>
    <li><a href="#user-content-Function_Convert">Convert</a></li>
<ul>
      <li><a href="#user-content-Function_Convert_bool">bool</a></li>
      <li><a href="#user-content-Function_Convert_float">float</a></li>
      <li><a href="#user-content-Function_Convert_int">int</a></li>
      <li><a href="#user-content-Function_Convert_string">string</a></li>
</ul>
    <li><a href="#user-content-Function_Date">Date</a></li>
<ul>
      <li><a href="#user-content-Function_Date_now">now</a></li>
      <li><a href="#user-content-Function_Date_parseTime">parseTime</a></li>
</ul>
    <li><a href="#user-content-Function_Encoding">Encoding</a></li>
<ul>
      <li><a href="#user-content-Function_Encoding_b64dec">b64dec</a></li>
      <li><a href="#user-content-Function_Encoding_b64enc">b64enc</a></li>
</ul>
    <li><a href="#user-content-Function_List">List</a></li>
<ul>
      <li><a href="#user-content-Function_List_first">first</a></li>
      <li><a href="#user-content-Function_List_includes">includes</a></li>
      <li><a href="#user-content-Function_List_last">last</a></li>
      <li><a href="#user-content-Function_List_list">list</a></li>
      <li><a href="#user-content-Function_List_map">map</a></li>
      <li><a href="#user-content-Function_List_reverse">reverse</a></li>
      <li><a href="#user-content-Function_List_sort">sort</a></li>
      <li><a href="#user-content-Function_List_uniq">uniq</a></li>
</ul>
    <li><a href="#user-content-Function_Math">Math</a></li>
<ul>
      <li><a href="#user-content-Function_Math_add">add</a></li>
      <li><a href="#user-content-Function_Math_ceil">ceil</a></li>
      <li><a href="#user-content-Function_Math_floor">floor</a></li>
      <li><a href="#user-content-Function_Math_max">max</a></li>
      <li><a href="#user-content-Function_Math_min">min</a></li>
      <li><a href="#user-content-Function_Math_mod">mod</a></li>
      <li><a href="#user-content-Function_Math_mul">mul</a></li>
      <li><a href="#user-content-Function_Math_quo">quo</a></li>
      <li><a href="#user-content-Function_Math_rem">rem</a></li>
      <li><a href="#user-content-Function_Math_round">round</a></li>
      <li><a href="#user-content-Function_Math_sub">sub</a></li>
</ul>
    <li><a href="#user-content-Function_Strings">Strings</a></li>
<ul>
      <li><a href="#user-content-Function_Strings_camelCase">camelCase</a></li>
      <li><a href="#user-content-Function_Strings_capitalize">capitalize</a></li>
      <li><a href="#user-content-Function_Strings_center">center</a></li>
      <li><a href="#user-content-Function_Strings_hasPrefix">hasPrefix</a></li>
      <li><a href="#user-content-Function_Strings_hasSuffix">hasSuffix</a></li>
      <li><a href="#user-content-Function_Strings_html">html</a></li>
      <li><a href="#user-content-Function_Strings_join">join</a></li>
      <li><a href="#user-content-Function_Strings_kebabCase">kebabCase</a></li>
      <li><a href="#user-content-Function_Strings_lower">lower</a></li>
      <li><a href="#user-content-Function_Strings_matchRegex">matchRegex</a></li>
      <li><a href="#user-content-Function_Strings_pascalCase">pascalCase</a></li>
      <li><a href="#user-content-Function_Strings_quote">quote</a></li>
      <li><a href="#user-content-Function_Strings_repeat">repeat</a></li>
      <li><a href="#user-content-Function_Strings_replace">replace</a></li>
      <li><a href="#user-content-Function_Strings_replaceN">replaceN</a></li>
      <li><a href="#user-content-Function_Strings_snakeCase">snakeCase</a></li>
      <li><a href="#user-content-Function_Strings_split">split</a></li>
      <li><a href="#user-content-Function_Strings_striptags">striptags</a></li>
      <li><a href="#user-content-Function_Strings_substr">substr</a></li>
      <li><a href="#user-content-Function_Strings_trim">trim</a></li>
      <li><a href="#user-content-Function_Strings_trimPrefix">trimPrefix</a></li>
      <li><a href="#user-content-Function_Strings_trimSuffix">trimSuffix</a></li>
      <li><a href="#user-content-Function_Strings_truncate">truncate</a></li>
      <li><a href="#user-content-Function_Strings_unquote">unquote</a></li>
      <li><a href="#user-content-Function_Strings_upper">upper</a></li>
      <li><a href="#user-content-Function_Strings_urlEscape">urlEscape</a></li>
      <li><a href="#user-content-Function_Strings_urlUnescape">urlUnescape</a></li>
      <li><a href="#user-content-Function_Strings_wordwrap">wordwrap</a></li>
</ul>
</ul>
</ul>

<h2><a id="user-content-Function" target="_self">Function</a></h2>
<h3><a id="user-content-Function__" target="_self">_</a></h3>

`_` is a no-op function that returns an empty string. It's useful to place a newline in the template. 
Example: 

```
{{- if .Ok}}
{{printf "ok: %v" .Ok}}
{{_}}
{{- end}}
```

<h3><a id="user-content-Function_Convert" target="_self">Convert</a></h3>
<h4><a id="user-content-Function_Convert_bool" target="_self">bool</a></h4>

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

<h4><a id="user-content-Function_Convert_float" target="_self">float</a></h4>

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

<h4><a id="user-content-Function_Convert_int" target="_self">int</a></h4>

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

<h4><a id="user-content-Function_Convert_string" target="_self">string</a></h4>

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

<h3><a id="user-content-Function_Date" target="_self">Date</a></h3>
<h4><a id="user-content-Function_Date_now" target="_self">now</a></h4>

`now` returns the current time. 
Example: 
```
{{now}}
```

Output: 
```
2024-09-12 15:04:05.999999999 +0000 UTC
```

<h4><a id="user-content-Function_Date_parseTime" target="_self">parseTime</a></h4>

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

<h3><a id="user-content-Function_Encoding" target="_self">Encoding</a></h3>
<h4><a id="user-content-Function_Encoding_b64dec" target="_self">b64dec</a></h4>

`b64dec` decodes a base64 encoded string. 
Example: 
```
{{b64dec "SGVsbG8sIFdvcmxkIQ=="}}
```

Output: 
```
Hello, World!
```

<h4><a id="user-content-Function_Encoding_b64enc" target="_self">b64enc</a></h4>

`b64enc` encodes a string to base64. 
Example: 
```
{{b64enc "Hello, World!"}}
```

Output: 
```
SGVsbG8sIFdvcmxkIQ==
```

<h3><a id="user-content-Function_List" target="_self">List</a></h3>
<h4><a id="user-content-Function_List_first" target="_self">first</a></h4>

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

<h4><a id="user-content-Function_List_includes" target="_self">includes</a></h4>

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

<h4><a id="user-content-Function_List_last" target="_self">last</a></h4>

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

<h4><a id="user-content-Function_List_list" target="_self">list</a></h4>

`list` creates a list from the given arguments. 
Example: 
```
{{list 1 2 3}}
```

Output: 
```
[1 2 3]
```

<h4><a id="user-content-Function_List_map" target="_self">map</a></h4>

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

<h4><a id="user-content-Function_List_reverse" target="_self">reverse</a></h4>

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

<h4><a id="user-content-Function_List_sort" target="_self">sort</a></h4>

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

<h4><a id="user-content-Function_List_uniq" target="_self">uniq</a></h4>

`uniq` removes duplicate elements from a list. 
Example: 
```
{{uniq (list 1 2 2 3 3 3)}}
```

Output: 
```
[1 2 3]
```

<h3><a id="user-content-Function_Math" target="_self">Math</a></h3>
<h4><a id="user-content-Function_Math_add" target="_self">add</a></h4>

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

<h4><a id="user-content-Function_Math_ceil" target="_self">ceil</a></h4>

`ceil` returns the least integer value greater than or equal to the input. 
Example: 
```
{{ceil 3.14}}
```

Output: 
```
4
```

<h4><a id="user-content-Function_Math_floor" target="_self">floor</a></h4>

`floor` returns the greatest integer value less than or equal to the input. 
Example: 
```
{{floor 3.14}}
```

Output: 
```
3
```

<h4><a id="user-content-Function_Math_max" target="_self">max</a></h4>

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

<h4><a id="user-content-Function_Math_min" target="_self">min</a></h4>

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

<h4><a id="user-content-Function_Math_mod" target="_self">mod</a></h4>

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

<h4><a id="user-content-Function_Math_mul" target="_self">mul</a></h4>

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

<h4><a id="user-content-Function_Math_quo" target="_self">quo</a></h4>

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

<h4><a id="user-content-Function_Math_rem" target="_self">rem</a></h4>

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

<h4><a id="user-content-Function_Math_round" target="_self">round</a></h4>

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

<h4><a id="user-content-Function_Math_sub" target="_self">sub</a></h4>

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

<h3><a id="user-content-Function_Strings" target="_self">Strings</a></h3>
<h4><a id="user-content-Function_Strings_camelCase" target="_self">camelCase</a></h4>

`camelCase` converts a string to camelCase. 
Example: 
```
{{camelCase "hello world"}}
```

Output: 
```
helloWorld
```

<h4><a id="user-content-Function_Strings_capitalize" target="_self">capitalize</a></h4>

`capitalize` capitalizes the first character of a string. 
Example: 
```
{{capitalize "hello"}}
```

Output: 
```
Hello
```

<h4><a id="user-content-Function_Strings_center" target="_self">center</a></h4>

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

<h4><a id="user-content-Function_Strings_hasPrefix" target="_self">hasPrefix</a></h4>

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

<h4><a id="user-content-Function_Strings_hasSuffix" target="_self">hasSuffix</a></h4>

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

<h4><a id="user-content-Function_Strings_html" target="_self">html</a></h4>

`html` escapes special characters in a string for use in HTML. 
Example: 
```
{{html "<script>alert('XSS')</script>"}}
```

Output: 
```
&lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;
```

<h4><a id="user-content-Function_Strings_join" target="_self">join</a></h4>

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

<h4><a id="user-content-Function_Strings_kebabCase" target="_self">kebabCase</a></h4>

`kebabCase` converts a string to kebab-case. 
Example: 
```
{{kebabCase "helloWorld"}}
```

Output: 
```
hello-world
```

<h4><a id="user-content-Function_Strings_lower" target="_self">lower</a></h4>

`lower` converts a string to lowercase. 
Example: 
```
{{lower "HELLO"}}
```

Output: 
```
hello
```

<h4><a id="user-content-Function_Strings_matchRegex" target="_self">matchRegex</a></h4>

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

<h4><a id="user-content-Function_Strings_pascalCase" target="_self">pascalCase</a></h4>

`pascalCase` converts a string to PascalCase. 
Example: 
```
{{pascalCase "hello world"}}
```

Output: 
```
HelloWorld
```

<h4><a id="user-content-Function_Strings_quote" target="_self">quote</a></h4>

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

<h4><a id="user-content-Function_Strings_repeat" target="_self">repeat</a></h4>

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

<h4><a id="user-content-Function_Strings_replace" target="_self">replace</a></h4>

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

<h4><a id="user-content-Function_Strings_replaceN" target="_self">replaceN</a></h4>

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

<h4><a id="user-content-Function_Strings_snakeCase" target="_self">snakeCase</a></h4>

`snakeCase` converts a string to snake_case. 
Example: 
```
{{snakeCase "helloWorld"}}
```

Output: 
```
hello_world
```

<h4><a id="user-content-Function_Strings_split" target="_self">split</a></h4>

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

<h4><a id="user-content-Function_Strings_striptags" target="_self">striptags</a></h4>

`striptags` removes HTML tags from a string. 
Example: 
```
{{striptags "<p>Hello <b>World</b>!</p>"}}
```

Output: 
```
Hello World!
```

<h4><a id="user-content-Function_Strings_substr" target="_self">substr</a></h4>

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

<h4><a id="user-content-Function_Strings_trim" target="_self">trim</a></h4>

`trim` removes leading and trailing whitespace from a string. 
Example: 
```
{{trim "  hello  "}}
```

Output: 
```
hello
```

<h4><a id="user-content-Function_Strings_trimPrefix" target="_self">trimPrefix</a></h4>

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

<h4><a id="user-content-Function_Strings_trimSuffix" target="_self">trimSuffix</a></h4>

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

<h4><a id="user-content-Function_Strings_truncate" target="_self">truncate</a></h4>

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

<h4><a id="user-content-Function_Strings_unquote" target="_self">unquote</a></h4>

`unquote` returns an unquoted string. 
Example: 
```
{{unquote "\"hello\""}}
```

Output: 
```
hello
```

<h4><a id="user-content-Function_Strings_upper" target="_self">upper</a></h4>

`upper` converts a string to uppercase. 
Example: 
```
{{upper "hello"}}
```

Output: 
```
HELLO
```

<h4><a id="user-content-Function_Strings_urlEscape" target="_self">urlEscape</a></h4>

`urlEscape` escapes a string for use in a URL query. 
Example: 
```
{{urlEscape "hello world"}}
```

Output: 
```
hello+world
```

<h4><a id="user-content-Function_Strings_urlUnescape" target="_self">urlUnescape</a></h4>

`urlUnescape` unescapes a URL query string. 
Example: 
```
{{urlUnescape "hello+world"}}
```

Output: 
```
hello world
```

<h4><a id="user-content-Function_Strings_wordwrap" target="_self">wordwrap</a></h4>

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

