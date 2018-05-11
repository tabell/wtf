---
title: "Textfile"
date: 2018-05-09T11:13:11-07:00
draft: false
---

## Description

Displays the contents of the specified text file in the widget.

<img src="/imgs/modules/textfile.png" width="320" height="133" alt="textfile screenshot" />

## Source Code

```bash
wtf/textfile/
```

## Required ENV Variables

None.

## Keyboard Commands

<span class="caption">Key:</span> `/` <br />
<span class="caption">Action:</span> Open/close the widget's help window.

<span class="caption">Key:</span> `o` <br />
<span class="caption">Action:</span> Opens the text file in whichever text editor is associated  with that file type.

## Configuration

```yaml
textfile:
  enabled: true
  filename: "notes.md"
  position:
    top: 5
    left: 4
    height: 2
    width: 1
  refreshInterval: 15
```

### Attributes

`enabled` <br />
Determines whether or not this module is executed and if its data displayed onscreen. <br />
Values: `true`, `false`.

`filename` <br />
The name of the file to be displayed in the widget. <br />
*Note:* Currently this file *must* reside in the `~/.wtf/` directory.
This is a <a href="https://github.com/senorprogrammer/wtf/issues/35">known bug</a>.

`position` <br />
Defines where in the grid this module's widget will be displayed. <br />

`refreshInterval` <br />
How often, in seconds, this module will update its data. <br />
Values: A positive integer, `0..n`.