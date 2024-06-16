# gh-green 
[![release](https://github.com/j-dunham/gh-green/actions/workflows/release.yml/badge.svg)](https://github.com/j-dunham/gh-green/actions/workflows/release.yml)

The `gh` entension that checks if you are "green" .i.e. have a green square in your contribution graph in GitHub!

## How to install
- Install the [GitHub CLI](https://cli.github.com)
- Install the extension `gh extension install j-dunham/gh-green`

## Try it out!
Run `gh green`


## Example outputs
for no contributions
```
*._.:*:._.:*:._.:*:._.:*:._.:*:._.:*:._.:*:._.:*
|                                              |
*  You haven't made any contributions... yet!  *
|                                              |
*._.:*:._.:*:._.:*:._.:*:._.:*:._.:*:._.:*:._.:*
```

with contributions
```
*._.:*:._.:*:._.:*:._.:*:._.:*
|                            |
*  You are green for today!  *
|                            |
*  totals:                   *
|  - 3 commits               |
*  - 1 issues                *
|  - 0 PRs                   |
*  - 2 PR reviews            *
|  - 0 repositories          |
*                            *
*._.:*:._.:*:._.:*:._.:*:._.:*
```