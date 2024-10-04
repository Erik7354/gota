# KIM & Nota (Encoder)

Nota is a basic serialization format developed by [Douglas Crockford](https://www.crockford.com) the creator of JSON.
In his own words from the WeAreDevelopers Conference 2024 - Nota is not meant to replace JSON but rather present an alternative that is as simple but more performant than JSON.

Nota uses its own string encoding [KIM](https://www.crockford.com/kim.html), instead of UTF-8. 
Nota also has its own format for floating point numbers called [DEC64](https://www.crockford.com/dec64.html).
Converting an IEEE754 float to DEC64 representation is surprisingly difficult to perform lossless.
[Here](https://www.reddit.com/r/programming/comments/28r8xt/dec64_is_intended_to_be_the_only_number_type_in/) is an interesting thread about DEC64.

Nota's design is based on the concept of a continuation bit at the beginning of every byte.
This would prevent a decoder of skipping over most datatypes like text, arrays and records without parsing them.

Given the complexity of the formats the KIM encoder and decoder is faster than the ones from the stdlib for UTF-8.
Benchmarks can be found in the kim package.

## Maintenance

This was a weekend project to dive deeper into KIM & Nota after hearing of it at the WAD24 conference.
The code is not maintained nor production ready and currently also lacks features like struct and float support and a Nota decoder.
