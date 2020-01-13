## Pixiq Architecture and design guidelines

Despite the main goals listed in [README](../README.md#project-goals) 
Pixiq has additional goals, mainly those of an architectural nature. 

Pixiq users expect:

+ no bugs
+ stable API
+ flexibility in replacing selected elements with others or adding new ones

Pixiq creators can fulfill those goals by:

+ writing automated tests (bearing in mind that 95% of tests should be unit tests)
+ designing the API carefully
  + making proof-of-concepts
  + using Test-Driven Development
  + using architectural patterns, SOLID principles and above all - hexagonal
  architecture
  + looking at competitive solutions
  + discussing solutions with others
+ making small, **independent** packages - see:
  + [image](../image)
  + [keyboard](../keyboard)
  + [loop](../loop)
+ in most cases new features should be added to new packages (unless something
  was missing from the beginning)
+ Pixiq should be more like a **library** not a framework. It basically means, 
  that the library user decides how to setup an application and how everything
  should work together, not the other way around. Sometimes using a premade setup 
  function or some fixed pipeline may save a few keystrokes, but this kind
  of functionality should be only available in addition to the usual means.
+ when Pixiq will reach the `v1.0.0` version it won't be possible to change the API
  (both syntactically and semantically) without introducing a new major version.
  Despite having a new git tag `v2.0.0`, new package `v2` will also be created.
+ therefore, at some point in time there will be a need to split the project into 
  pieces in order to support different versioning for each module. Stable modules
  (such as `image` or `keyboard`) will have `v1.x.x` version, unstable ones
  (such as `devtool`) will stay at `v0.x.x`. Thanks to that Pixiq developers 
  will still be able to introduce introduce incompatible changes to unstable modules.
  Such freedom greatly improves creativity.

### Reasoning behind design decisions (so far)

Why there is no abstraction for opening windows?

> Because it is really hard to design such abstraction. There are way too many
platforms varying in possibilities (PCs, Macs, mobile devices etc.). But we are
open for [proposals](https://github.com/jacekolszak/pixiq/issues). There is 
a chance that it will be feasible to create such abstraction just for PCs 
(Win, Mac and Linux).


The [opengl](../opengl) package uses [GLFW](https://www.glfw.org/), but for 
some reason it does not provide all the features of the mentioned library,
for example it does not allow to change the look of the cursor or to set
the window transparency. Why?

> Because we haven't had time to do it yet. If something is really important
for you then propably it is time to submit an  [Issue](https://github.com/jacekolszak/pixiq/issues) 
or maybe even make a change by submitting a [Pull Request](https://github.com/jacekolszak/pixiq/pulls).
The `opengl.OpenGL` struct does not need to implement any abstractions
for opening windows, therefore it might be extended freely.


Why you are still using OpenGL in 2020?

> Because the API is easy to use (compared to Vulkan for example), is multiplatform (yes, it still works on MacOS) and has decent support in Go.

