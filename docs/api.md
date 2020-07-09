# APIs for AIArts

## 接口风格

### 1. URL
* 所有api以`/ai_arts/api`开头
* 不同模块url prefix为`/ai_arts/api/<modulename>`

### 2. 返回值
* 返回值为json格式
* 返回值分为三部分:
    - code: 错误码，integer，0为正确，其他为5位数字，具体错误见错误码文档(之后补充)
    - msg: 接口正确为`success`，错误为错误信息
    - data: 返回的具体数据
* 各接口返回值后续补充

### 3. 传参
* GET请求，参数写在query_string
* POST请求，参数使用json body
* 上传文件，使用form表单提交
* 各接口参数待补充

## 接口列表

### 1. 代码开发 `/ai_arts/api/codes`
实现上，调用平台job，主要使用jupyter的endpoint。

#### 1.1 `GET /ai_arts/api/codes` 获取代码开发列表

#### 1.2 `POST /ai_arts/api/codes` 创建新的代码开发

#### 1.3 `GET /ai_arts/api/codes/:id` 获取单个代码开发信息

#### 1.4 `GET /ai_arts/api/codes/:id/endpoint` 跳转到该id的jupyter url(也可前端直接跳转)

#### 1.5 `DELETE /ai_arts/api/codes/:id` 删除单个代码开发

#### 1.6 `POST /ai_arts/api/codes/:id/description` 更新描述信息

### 2. 模型管理 `/ai_arts/api/models`

#### 2.1 `GET /ai_arts/api/models` 获取模型列表

#### 2.2 `POST /ai_arts/api/models/upload` 上传模型

#### 2.3 `POST /ai_arts/api/models` 创建模型

#### 2.4 `GET /ai_arts/api/models/:id/download` 下载模型

#### 2.5 `DELETE /ai_arts/api/models/:id` 删除模型

#### 2.6 模型注册？类似数据集注册

### 3. 模型训练 `/ai_arts/api/trainings`

#### 3.1 `GET /ai_arts/api/trainings` 获取训练列表

#### 3.2 `GET /ai_arts/api/trainings/:id` 获取训练详情

#### 3.3 `GET /ai_arts/api/trainings/:id/log` 训练日志获取？

#### 3.4 `POST /ai_arts/api/trainings` 创建训练

#### 3.4 `POST /ai_arts/api/trainings/:id/stop` 停止训练

#### 3.5 `DELETE /ai_arts/api/trainings/:id` 删除训练

### 4. 数据集管理 `/ai_arts/api/datasets`

#### 4.1 `GET /ai_arts/api/datasets` 获取数据集列表

#### 4.2 `GET /ai_arts/api/datasets/:id` 获取数据集详情

#### 4.3 `POST /ai_arts/api/datasets` 创建数据集

#### 4.4 `POST /ai_arts/api/datasets/upload` 上传数据集

#### 4.5 `POST /ai_arts/api/datasets/register` 注册数据集

#### 4.6 `GET /ai_arts/api/datasets/:id/download` 下载数据集

#### 4.7 `DELETE /ai_arts/api/datasets/:id` 删除数据集

### 5. 数据标注 `/ai_arts/api/annotations`
此处跳转到其他项目

### 6. 推理服务 `/ai_arts/api/inferences`

#### 6.1 `GET /ai_arts/api/inferences` 获取推理列表

#### 6.2 `GET /ai_arts/api/inferences/:id` 获取推理详情

#### 6.3 `POST /ai_arts/api/inferences` 创建推理

#### 6.4 `POST /ai_arts/api/inferences/:id/upload_image` 上传图片

#### 6.5 `POST /ai_arts/api/inferences/:id/recognition` 图片识别
