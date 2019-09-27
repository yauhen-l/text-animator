This library provides only one method `DrawGif` which generates animated GIF image, based on an imput.

Here is an example:
```
f, _ := os.OpenFile("out.gif", os.O_WRONLY|os.O_CREATE, 0600)
defer f.Close()

err := DrawGif(DefaultFace(), []string{
"ba ",
"   DUM!",
" Tss",
}, []int{10, 50, 100}, f)

```

`out.gif` file in this case will look like this:

![ba-DUM-Tss](https://media.giphy.com/media/eIrsVaZIYHaeX1czoh/giphy.gif)


This library I use in telegram bot: https://telegram.me/text_shots_bot

