# yc-log-ui
Alternative Yandex Cloud Logging UI

## Установка 
Подключение Helm репозитория
```shell
helm repo add yc-log-ui https://constantable.github.io/yc-log-ui/
```
Установка
```shell
helm upgrade --install yc-log-ui .helm -n kube-logging --values my-values.yaml
```

## Настройки

| Переменная                            | Значение                          |
|---------------------------------------|-----------------------------------|
| backend.parameters.service_account_id | Id сервисного аккаунта            |
| backend.parameters.log_group_id       | Id лог-группы                     |
| backend.parameters.public_key         | Публичный ключ (base64 encoded)   |
| backend.parameters.private_key        | Приватный ключ (base64 encoded)   |
| backend.parameters.key_id             | Id ключа                          |
| ingress.hosts[]                       | Настройки хостов ingress          |