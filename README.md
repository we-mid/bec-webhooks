Todo
- [ ] testing
- [x] stash and pop by specific name
- [x] detect if ws is clean
- [x] handle conflicts (leave it for manually resolving)

Article: [GitHub Webhooks 技术方案及落地](https://bec.today/fx/?draft/github-webhooks)

```txt
### Demo ###
15|bec-web | 12-06 20:19: [receiving] push: repo=user/foo, sender=fritx
15|bec-web | 12-06 20:19: [done] pull: path="/path/to/foo"
15|bec-web | 12-06 22:34: [receiving] push: repo=user/foo, sender=fritx
15|bec-web | 12-06 22:34: [done] pull: path="/path/to/foo"
15|bec-web | 12-06 23:14: [receiving] push: repo=org/bar, sender=fritx
15|bec-web | 12-06 23:14: [done] pull: path="/path/to/bar"
```
