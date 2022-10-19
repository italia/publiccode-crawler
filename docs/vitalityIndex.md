<!--
    This file is linked from https://developers.italia.it,
    DON'T REMOVE/RENAME it without updating the link first.
-->

## Vitality Index

This document shows how the vitality index is calculated.

The vitality index represents how much a certain repository has been active in the
last period of time.
As such, it focuses on the following data:

* Code activity: the number of commit and merges on a daily basis;
* Release history: the number of daily releases;
* User community: the number of unique authors;
* Longevity: the age of the project.

Right now the algorithm is using a window of the last 60 days.

### User community

This indicator represents the number of users that authored a commit in the last
days.
Knowing how many users interacted with the code provides an indication of the
community around the code. Having an active community means that several users
interacted with the codebase in the last months.
The user community indicator has a minimum value of 4 points and a maximum of 36.
Please check
[this](https://github.com/italia/publiccode-crawler/blob/663c661ca3b0d6e1578f24c7be97fd35e28abe87/crawler/vitality-ranges.yml#L1-L29)
file to see how such range is divided.

### Code Activity

This indicator represents the number of activities performed in the last amount
of days. As such, two actions are considered:

1. the number of commits;
2. the number of merges.

The code activity indicator has a minimum value of 2 and a maximum of 60.
Please check
[this](https://github.com/italia/publiccode-crawler/blob/663c661ca3b0d6e1578f24c7be97fd35e28abe87/crawler/vitality-ranges.yml#L31-L62)
file to see the details of the indicator.

### Release History

This indicator tells how many releases have been done in the last period of time.
From a git repository, a quick way to tell this is by analyzing the `tags`.
This indicator inspects the number of tags and provides a minimum of 20 points,
when there have been from 0 to 1 releases, to a maximum of 50 when there have
been from 4 to 100 releases. Please check
[this](https://github.com/italia/publiccode-crawler/blob/663c661ca3b0d6e1578f24c7be97fd35e28abe87/crawler/vitality-ranges.yml#L64-L77)
reference file for the details.

### Longevity

The longevity is basically the repository age. As such, it is calculated by
extracting the date of the oldest commit and by calculating the different
between now and then. From 0 to 1 year time means 20 points whilst from 2 years
on means the maximum, that is 35 points.
[Here](https://github.com/italia/publiccode-crawler/blob/663c661ca3b0d6e1578f24c7be97fd35e28abe87/crawler/vitality-ranges.yml#L79-L89)
you can see more about this.

### Final calculation

Since all the above indicators are calculated on a daily basis the final
vitality index is simply a sum of the average of the 4 categories with some
final corrections in case it overflows 100%.
See
[these](https://github.com/italia/publiccode-crawler/blob/663c661ca3b0d6e1578f24c7be97fd35e28abe87/crawler/crawler/repo_activity.go#L100-L117)
lines for more details about it.
