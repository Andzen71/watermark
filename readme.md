# Readme

Сервис сделан для того, чтобы ему отправляли каринку, а получили картинку с водяным знаком.

Сервер запускается на 8080, поход на /watermark пост запросом с передачей в теле запроса картинки вызывает обработку запроса и в ответе должна прийти картинка с водяным знаком.

Команда для отправки картинки на сервер через форму (я так не хочу)

```bash
curl -X POST -F "image=@perf1.png" http://localhost:8080/watermark > watermarked_tree.png
```

Команда для отправки картинки на сервер через тело (как я хочу)

```bash
 curl localhost:8081/watermark -X POST --data-binary "@perf1.png" --output file.png -H "Content-type: image/png"
```

Заголовки нужны потому, что без них курл отправляет картинку с таким заголовком `Content-Type: application/octet-stream`

❌ Закончил на том, что кажется, что курл тело запроса оборачивает в base64 и поэтому код не может картинку заэнкодить.
Также у меня есть код, который энкодит эту же картинку, но берёт её с диска.

✅ На самом деле у curl есть параметр -d, что является алиасом к --data-ascii - а эта штука все символы, которые видит - url-encoded'ит