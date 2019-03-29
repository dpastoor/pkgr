# pkgr

[![asciicast](https://asciinema.org/a/wgcPBvCMtEwhpdW793MBjgSi2.svg)](https://asciinema.org/a/wgcPBvCMtEwhpdW793MBjgSi2)

# THIS IS CURRENTLY A WIP, however is getting close for user testing. Check back soon for more comprehensive user docs

# What is pkgr?

`pkgr` is a rethinking of the way packages are managed in R. Namely, it embraces
the declarative philosophy of defining _ideal state_ of the entire system, and working
towards achieving that objective. Furthermore, `pkgr` is built with a focus on reproducibility
and auditability of what is going on, a vital component for the pharmaceutical sciences + enterprises.

# Why pkgr?

`install.packages` and friends such as `remotes::install_github` have a subtle weakness --
they are not good at controlling desired global state. There are some knobs that
can be turned, but overall their APIs are generally not what the user _actually_ needs. Rather, they
are the mechanism by which the user can strive towards their needs, in a forceably iterative fashion.

With pkgr, you can, in a **parallel-processed** manner, do things like:
- Install a number of packages from various repositories, when specific packages must be pulled from specific repositories
- Install `Suggested` packages only for a subset of all packages you'd like to install
- Customize the installation behavior of a single package in a documentable and reproducible way
  - Set custom Makevars for a package that persist across system installations
  - Install source versions of some packages but binaries for others
- **Understand how your R environment will be changed _before_ performing an installation or action.**

Today, packages are highly interwoven. Best practices have pushed towards small, well-scoped packages that
do behaviors well. For example, rather than just having plyr, we now use dplyr+purrr to achieve
the same set of responsibilities (dealing with dataframes + dealing with other list/vector objects in an iterative way).
As such, it is becoming increasingly difficult to manage the _set_ of packages in a transparent and robust
way.

# How it works

`pkgr` is a command line utility with several top level commands. The two primary commands are:

```bash
pkgr plan # show what would happen if install is run
pkgr install # install the packages specified in pkgr.config
```
The actions are controlled by a configuration file that specifies the desired global state, namely,
by defining the top level packages a user cares about, as well as specific configuration customizations.

For example, a pkgr configuration file might look like:

```yaml
Version: 1
# top level packages
Packages:
  - rmarkdown
  - bitops
  - caTools
  - knitr
  - tidyverse
  - shiny
  - logrrr

# any repositories, order matters
Repos:
  - gh_dev: "https://metrumresearchgroup.github.io/rpkgs/gh_dev"
  - CRAN: "https://cran.microsoft.com/snapshot/2018-11-18"

# path to install packages to
Library: "path/to/install/library"

# package specific customizations
Customizations:
  - tidyverse:
      Suggests: true
```

When you run `pkgr install` with this as your _pkgr.config_ file, pkgr will download and
install the packages rmarkdown, bitops, calToools, knitr, tidyverse, shiny, logrrr,
and any dependencies that those packages require. Since the "gh_dev" repository is listed first,
pkgr will search "gh_dev" for those packages before it looks to "CRAN".

If you want to see everything that pkgr is going to install before actually installing, simply run `pkgr plan` and take a look.



How about a more complex example?

Let's say you're working on an OSX machine.
On CRAN, for OSX, the package `devtools` (v2.x) is currently available as source,
but the binary is still v1.13. You want the latest version of devtools, so you'll need to build it from source.
However, you still want to install from binaries (the default behavior for OSX) for everything else in your environment.
With pkgr, you can set a `Customization` for `devtools` using `Type: source`

```yaml
Version: 1
# top level packages
Packages:
  - rmarkdown
  - shiny
  - devtools

# any repositories, order matters
Repos:
   - CRAN: "https://cran.microsoft.com/snapshot/2018-11-18"

Library: "path/to/install/library"

# can cache both the source and installed binary versions of packages
Cache: "path/to/global/cache"

# can log the actions and outcomes to a file for debugging and auditing
Logging:
  File: pkgr-install.log

Customizations:
  - devtools:
      Type: source
```

With this customization in your config file, pkgr will install from sources for devtools.
For everything else, the default install behavior will stay in effect.

For a third example, here is a configuration that also pulls from bioconductor:

```
Version: 1
# top level packages
Packages:
  - magrittr
  - rlang
  - ggplot2
  - dplyr
  - tidyr
  - plotly
  - VennDiagram
  - aws.s3
  - data.table
  - forcats
  - preprocessCore
  - loomR
  - ggthemes
  - reshape

# any repositories, order matters
Repos:
  - CRAN: "https://cran.microsoft.com/snapshot/2018-11-18"
  - BioCsoft: "https://bioconductor.org/packages/3.8/bioc"
  - BioCann: "https://bioconductor.org/packages/3.8/data/annotation"
  - BioCexp: "https://bioconductor.org/packages/3.8/data/experiment"
  - BioCworkflows: "https://bioconductor.org/packages/3.8/workflows"

# path to install packages to
Library: pkgs

Cache: pkgcache
Logging:
  File: pkgr-install.log
```

## Pkgr and [Packrat](https://rstudio.github.io/packrat/)

**Pkgr is not a replacement for Packrat**. Packrat is a tool to capture the state
of your R environment and isolate it from outside modification.
Where Packrat often falls short, however, is in the restoration said environment.
Packrat::restore() restores packages in an iterative fashion, which is a
time-consuming process that doesn't always play nice with packages hosted outside
of CRAN (such as packages hosted on GitHub). Additionally, packrat::restore()
provides the user with almost no control of the _order_ in which packages are
restored -- it always installs packages alphabetically by name.

Pkgr solves these issues by:
  - Installing packages quickly in parallelized layers (determined by the dependency tree)
  - Allowing users to control things like what repo a given package is retrieved from
  - Showing users a holistic view of their R Environment (`pkgr inspect --deps --tree`) and how that environment would be changed on another install (`pkgr plan`)

## More info to come as we progress!

As we continue development, we intend to answer the questions:
- How does pkgr integrate with [Packrat](https://rstudio.github.io/packrat/)?
- How does pkgr integrate with [RStudio Package Manager](https://www.rstudio.com/products/package-manager/)?
- How much faster is pkgr than its peers (such as install.packages)?
- What are some of the downsides of pkgr?


## API options

package declaration can become nuanced as the user desires to customize
specifically where a package is pulled from.
Given a set of repositories, the default R tooling
will stop after the package is found and use that.
In some cases the user may prefer to explicitly
declare which repository a package may come from, especially when pulling from
an environment where multiple repos are specified.

The remotes API provides a rich experience around customizing the installation
behavior from external repositories such as github, where combinations such as

- tidyverse/dplyr#12345 - install from PR 123
- github::tidyverse/dplyr#12345 - install from github

The intent of this tooling (for now) is to provide a more modular experience,
in which the ways packages can be identified is minimized, and upstream tooling
can coalesce packages from many of the scenarios outlined.

As such, the biggest focus is on targetting packages placed in a specific _repository_.
With this in mind, the question remains whether the remotes API should be followed.
The concern is that the repo::package pattern slightly obfusicates the package.
This is less noticeable when the package is previously declared in the Imports/Depends
statement of a DESCRIPTION file, however when packages become the forefront
of the requirements.

Some of the potential API designs are:

```yaml
Packages:
  - repository::package
  - package@repository
  - package
    repo: repository
```

### Package dependencies:

By default, packages will need Imports/Depends/LinkingTo to make
sure the packages can work successfully.

```yaml
Packages:
  - PKPDmisc
    Suggests: true
  - dplyr
```

The benefit is customizations related to package requirements
are immediately visible. The downside is it "pollutes" the
packages list.

```yaml
Packages:
  - PKPDmisc
  - dplyr

Customizations:
  PKPDmisc:
    Suggests: true
```


## Assumptions (for now)

Making this tool bulletproof will take a significant effort over time. To bring confidence for use day-to-day
we must clearly outline the assumptions we are making internally to provide guidance on the areas
this tool may not be suitable for, or to explain unexpected behavior.

* Package/versions from a given repo will not change over time
  * if pkgx_0.1.0 is downloaded from repoY, we do not need to check each time that pkgx is consistent
  * this allows simple caching without doing hash comparisons (for now)

R package management

## Install Strategy Background

One of the problems with the full layered implementation is the longest install time dictates the entire layer
installation install. Originally, we did not know if this would be a huge problem, however it was quickly
evident that this was not the case.

For example, when look at the installation layers given a request for ggplot2, the following was
the installation timing. For layer 1, the second longest install time was Rcpp (39 seconds), with most
other packages coming in less than 10 seconds.

| layer  |package | duration|
|:-------|:-------|--------:|
|   1    |stringi |   159.39|
|   2    |Matrix  |    69.98|
|   3    |mgcv    |    34.45|
|   4    |tibble  |     2.37|
|   5    |ggplot2 |    12.12|

Furthermore, when looking at subsequent layers, neither Matrix or mgcv and its dependencies have any relation
to stringi, so there is no reason to wait for the layer to complete.

# Development

run all tests with tabular output:

```
go test ./... -json -cover | tparse -all
```
