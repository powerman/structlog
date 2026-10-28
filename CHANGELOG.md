# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.10.0] - 2026-10-28

### 🔔 Changed

- Require Go 1.26 by @powerman in [b3db6bd]

### 📦️ Dependencies

- **(deps)** Upgrade go and golangci-lint by @powerman in [c9f9ce0]
- **(deps)** Upgrade go modules by @powerman in [b72392f]
- **(deps)** Upgrade go and golangci-lint by @powerman in [7be98fc]
- **(deps)** Mise up by @powerman in [ada66d0]

[0.10.0]: https://github.com/powerman/structlog/compare/v0.9.0..v0.10.0
[c9f9ce0]: https://github.com/powerman/structlog/commit/c9f9ce0523fec0df5f33d9762dc02ddd4336ea66
[b72392f]: https://github.com/powerman/structlog/commit/b72392fcf1ac6bb224b9782ca97e795e1f5fb72b
[7be98fc]: https://github.com/powerman/structlog/commit/7be98fc05a75b732cf42cc72441ed857fdeda51e
[b3db6bd]: https://github.com/powerman/structlog/commit/b3db6bd1464a02a9e492833cabb3a053400b0f03
[ada66d0]: https://github.com/powerman/structlog/commit/ada66d0a49989b87b3d8efe2d8ab254b25709c8f

## [0.9.0] - 2026-07-12

### 📦️ Dependencies

- **(deps)** Upgrade github.com/powerman/check to v1.13.0 by @powerman in [9cd760f]

[0.9.0]: https://github.com/powerman/structlog/compare/v0.8.1..v0.9.0
[9cd760f]: https://github.com/powerman/structlog/commit/9cd760fa31db37fe4ce27d52ca2a4a1174de68ad

## [0.8.1] - 2026-04-07

### 📦️ Dependencies

- **(deps)** Upgrade by @powerman in [ffde093]

[0.8.1]: https://github.com/powerman/structlog/compare/v0.8.0..v0.8.1
[ffde093]: https://github.com/powerman/structlog/commit/ffde093de24611e267807821392c6cfc51027609

## [0.8.0] - 2025-08-28

### 📦️ Dependencies

- **(deps)** Bump google.golang.org/protobuf from 1.27.0 to 1.33.0 by @dependabot[bot] in [#29]
- **(deps)** Bump google.golang.org/grpc from 1.38.0 to 1.56.3 by @dependabot[bot] in [#31]
- **(deps)** Bump `github.com/prometheus/client_golang` by @dependabot[bot] in [#33]
- **(deps)** Bump gopkg.in/yaml.v3 by @dependabot[bot] in [#34]

[0.8.0]: https://github.com/powerman/structlog/compare/v0.7.3..v0.8.0
[#29]: https://github.com/powerman/structlog/pull/29
[#31]: https://github.com/powerman/structlog/pull/31
[#33]: https://github.com/powerman/structlog/pull/33
[#34]: https://github.com/powerman/structlog/pull/34

## [0.7.3] - 2021-08-12

### 🐛 Fixed

- Avoid json marshal errors by @powerman in [#18]

[0.7.3]: https://github.com/powerman/structlog/compare/v0.7.2..v0.7.3
[#18]: https://github.com/powerman/structlog/pull/18

## [0.7.2] - 2021-07-12

### 🐛 Fixed

- Marshal msg as string in JSON by @powerman in [#17]

[0.7.2]: https://github.com/powerman/structlog/compare/v0.7.1..v0.7.2
[#17]: https://github.com/powerman/structlog/pull/17

## [0.7.1] - 2020-05-12

### 🔔 Changed

- Include KeyTime in default prefix by @powerman in [#10]

[0.7.1]: https://github.com/powerman/structlog/compare/v0.7.0..v0.7.1
[#10]: https://github.com/powerman/structlog/pull/10

## [0.7.0] - 2020-05-12

### 🚀 Added

- Support KeyTime=Auto for text output by @powerman in [#9]

### 🔔 Changed

- Convert README to markdown by @powerman in [6da519d]

[0.7.0]: https://github.com/powerman/structlog/compare/v0.6.0..v0.7.0
[6da519d]: https://github.com/powerman/structlog/commit/6da519ddfb81a81cce8293a511ac177352e99bd8
[#9]: https://github.com/powerman/structlog/pull/9

## [0.6.0] - 2019-10-30

### 🚀 Added

- Add support for errors.Unwrap by @powerman in [4e8b822]

### 📦️ Dependencies

- **(deps)** Bump github.com/powerman/check from 1.0.1 to 1.1.0 by @dependabot-preview[bot] in [#4]

[0.6.0]: https://github.com/powerman/structlog/compare/v0.5.0..v0.6.0
[#4]: https://github.com/powerman/structlog/pull/4
[4e8b822]: https://github.com/powerman/structlog/commit/4e8b8226d8ce8197f3f1207f238e5264591827e9

## [0.5.0] - 2019-05-30

### 🔔 Changed

- Add WrapErr example by @powerman in [7b46dc8]
- WrapErr return nil on nil error by @powerman in [ddb92f5]

[0.5.0]: https://github.com/powerman/structlog/compare/v0.4.0..v0.5.0
[7b46dc8]: https://github.com/powerman/structlog/commit/7b46dc888173f7a47c87e8de568f422d87d0267b
[ddb92f5]: https://github.com/powerman/structlog/commit/ddb92f5763e3e32ac3ea2d440c2422be39d6332d

## [0.4.0] - 2019-05-26

### 🚀 Added

- Add WrapErr by @powerman in [4556010]

### 🔔 Changed

- Add package overview by @powerman in [893e4f9]
- Cleanup by @powerman in [9ef6dfa]

[0.4.0]: https://github.com/powerman/structlog/compare/v0.3.0..v0.4.0
[893e4f9]: https://github.com/powerman/structlog/commit/893e4f971307e53fb197671ed708114d58abb864
[9ef6dfa]: https://github.com/powerman/structlog/commit/9ef6dfa52f3f10fe7c2dcf35993dfcf1b65fc1b8
[4556010]: https://github.com/powerman/structlog/commit/45560106815aba7a2272a3a8258b202925b7cfb3

## [0.3.0] - 2019-04-14

### 🚀 Added

- Add NewContext/FromContext by @powerman in [4118f63]

[0.3.0]: https://github.com/powerman/structlog/compare/v0.2.0..v0.3.0
[4118f63]: https://github.com/powerman/structlog/commit/4118f638248ae1485d669a6db774094b55daf383

## [0.2.0] - 2019-03-04

### 🚀 Added

- Add SetPrinter by @powerman in [31d5dc0]
- Add SetOutput by @powerman in [f76c909]

### 🔔 Changed

- Update README by @powerman in [faaf8db]

### 🐛 Fixed

- Source file/line of panic in Recover by @powerman in [a91ea57]

[0.2.0]: https://github.com/powerman/structlog/compare/v0.1.1..v0.2.0
[31d5dc0]: https://github.com/powerman/structlog/commit/31d5dc0df92c89b9313f712dbb06f2027f44020b
[f76c909]: https://github.com/powerman/structlog/commit/f76c90997b5bf410eac4e2b88bdabc02d0629cab
[a91ea57]: https://github.com/powerman/structlog/commit/a91ea5742bd5a496f343755dc827d940e2da11bb
[faaf8db]: https://github.com/powerman/structlog/commit/faaf8dbad08282642542fff5fa7e0beae836d7ce

## [0.1.1] - 2018-11-15

### 🔔 Changed

- Add go.mod by @powerman in [619f80e]
- Setup linter by @powerman in [fe22a6d]
- Switch to CircleCI by @powerman in [a8de34e]

[0.1.1]: https://github.com/powerman/structlog/compare/%40%7B10year%7D..v0.1.1
[619f80e]: https://github.com/powerman/structlog/commit/619f80e6e315bc74dcdb7f72aadf8c961c2be282
[fe22a6d]: https://github.com/powerman/structlog/commit/fe22a6d9d52cbe4470152bd2e482ea262bafaeae
[a8de34e]: https://github.com/powerman/structlog/commit/a8de34e542945a91e15098fcd361947d6521da76

<!-- generated by git-cliff -->
