# Lessons Learned

---
### Request Handling
Functions used:
- http.ServeMux : for request multiplexing
- http.HandlerFunc : converts a normal function into a Handler interface
- http.ListenAndServe : listens on the specified addr and then calls the given Handler for request handling
- http.Redirect : replies to the request with a redirect to a url
---
### YAML
- Recursive acronym for YAML Aint Markup Language
- golang provides a package to decode yaml file into an output value which can be a map or pointers
- A struct can be declared according to the specified YAML file to capture the YAML data
