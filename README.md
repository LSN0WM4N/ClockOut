# ClockOut

Planned to be an OpenSource solution for small business that have no to pay 
for a super massive software in order to control clock in/clock out.

### Project Status
Basic Core development

Next steps...

Create the handlers and listeners for system events.
I have in mind an http based connection restricted to localhost.

The main plan is separate the project into pieces that make their own 
function without knowing what the other piece do.

```
[client] <---> [listener <--channels--> dispatcher] <---> [core]
```

Update: 
I have no more time for today :(

```
TODOS:
  Add endpoints for CRUD employees
  Make the dispatcher able to accept more handlers
  Think a little bit more about the organization (is getting a little messy)
  
  Drink more water
```

### Principal Features in mind

* [ ] Manage times
* [ ] Generate statistics 
* [ ] Robust Backend with an easy frontend available
* [ ] Create extensions to work not only with software (Contactless cards is my primary goal)

### Contrib
If you want to help me or share any opinion you can make a PR with a full description of what
and why are you doing it.
I'm very happy to get feedback at any time.

Feel free also to text me to my [email](mailto:rodriguezrondafrancis@gmail.com) 