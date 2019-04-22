## Requirements

The list of redirection should be maintained in a command line tool, what can:

- [ ] Manipulate YAML config file. Where the redirection list peristently stored.
- [ ] Implement append to the list: urlshorten configure -a dogs -u www.dogs.com
- [ ] Implement remove from the list: urlshorten -d dogs
- [x] List redirections: urlshorten -l
- [ ] Run HTTP server on a given port: urlshorten run -p 8080
- [ ] Prints usage info: urlshorten -h
