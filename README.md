# goTV

![build:?](https://travis-ci.org/prospero78/goTV.svg)

Пользовательский интерфейс командной строки **CUI** по аналогии с **TurboVision**) с поддержкой встроенных тем. Смотрите скриншоты рабочих примеров в конце этого **README**.

## Установка

```bash
    go get -u github.com/prospero78/goTV
```

## Текущая версия

Текущая версия 1.2.1. Для подробностей смотрите [изменения](./changelog.md).

## Приложения, использующие библиотеку

- [Terminal FB2 reader(termfb2)](https://github.com/VladimirMarkelov/termfb2)

## Документация

- [Введение](./docs/intro.md)
- [Перед стартом](./docs/hello.md)
- [Менеджер компоновки](./docs/layout.md)
- [Основные стандартные виджеты, их методы и свойства](./docs/widget.md)
- [Окна](./docs/window.md)
- [Предопределённые горячие клавиши](./docs/hotkeys.md)

## Список доступных виджетов

- `Window` (Главный контейнер виджетов - с максимизацией, порядком отрисовки и другими различными свойствами)
- `Label` (Horizontal и Vertical with basic color control tags)
- `Button` (Simple push button control)
- `EditField` (One line text edit control with basic clipboard control)
- `ListBox` (string list control with vertical scroll)
- `TextView` (ListBox-alike control with vertical and horizontal scroll, and wordwrap mode)
- `ProgressBar` (Vertical and horizontal. The latter one supports custom text over control)
- `Frame` (A decorative control that can be a container for other controls as well)
- Scrollable frame
- `CheckBox` (Simple check box)
- `Radio` (Simple radio button. Useless alone - should be used along with RadioGroup)
- `RadioGroup` (Non-visual control to manage a group of a few RadioButtons)
- `ConfirmationDialog` (modal View to ask a user confirmation, button titles are custom)
- `SelectDialog` (modal View to ask a user to select an item from the list - list can be ListBox or RadioGroup)
- `SelectEditDialog` (modal View to ask a user to enter a value)
- `BarChart` (Horizontal bar chart without scroll)
- `SparkChart` (Show tabular data as a bar graph)
- `GridView` (Table to show structured data - only virtual and readonly mode with scroll support)
- [FilePicker](./docs/fselect.md)
- `LoginDialog` - a simple authorization dialog with two fields: Username and Password
- `TextDisplay` - a "virtual" text view control: it does not store any data, every time it needs to draw its line it requests the line from external source by line ID

## Скриншоты

The main demo (theme changing and radio group control)

![Главное демо](./demos/clui_demo_main.gif)

The screencast of demo:

![Демо библиотеки](./demos/demo.gif)

The library is in the very beginning but it can be used to create working utilities: below is the example of my Dilbert comix downloader:

![Закачка Дилберта](./demos/dilbert_demo.gif)

## Лицензия

**goTV** распространяется под лицензией [BSD-2-Clause](./docs/LICENSE.md).
