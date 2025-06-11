"use server"
import { Readable } from 'stream';
import { promises as fs } from 'fs';
import path from 'path';
const { v4: uuidv4 } = require('uuid');

// 生成随机文件名
function generateRandomFileName(fileName) {
  const randomKey = uuidv4(); 
  const parts = fileName.split(".");
  const extension = parts.length > 1 ? `.${parts.pop()}` : '';
  return `${randomKey}${extension}`; // 拼接随机字符串和文件扩展名
}

// 上传文件的函数
export default async(data)=> {
  const file = data.get("file")
  if (!file) {
    throw new Error("No file uploaded")
  }
  const newFilename = generateRandomFileName(file.name);
  // 读取文件内容
  // 将 File 对象转换为 Node.js 的 ReadableStream
  const readableStream = Readable.fromWeb(file.stream());
  // 读取文件内容到 Buffer
  const chunks = [];
  for await (const chunk of readableStream) {
    chunks.push(chunk);
  }
  const buffer = Buffer.concat(chunks);

  // 保存文件到本地
  const uploadDir = path.join(process.cwd(), 'public/images');
  console.log('本地路径：', uploadDir);
  await fs.mkdir(uploadDir, { recursive: true });
  const filePath = path.join(uploadDir, newFilename);
  await fs.writeFile(filePath, buffer);

  // 返回文件的本地访问路径
  const fileUrl = `/images/${newFilename}`;
  console.log('文件保存路径：', fileUrl);
  return fileUrl;
}
// "use server"
// import { Readable } from 'stream';
// const qiniu = require('qiniu');
// const { v4: uuidv4 } = require('uuid');

// const accessKey = 'mxf_OypOQR4g_RDujwvw7G9q0Jo3X8mpvKiBvVje';
// const secretKey = 'vJ4zoiPKs-qCTiVe9FjKYXKv3Oo8jyQF4LHWGLH9';
// const domainOfBucket= 'ssmqli488.hn-bkt.clouddn.com';
// const mac = new qiniu.auth.digest.Mac(accessKey, secretKey);
// // 设置要操作的空间名
// const bucket = 'bitdancepanic';

// // 生成上传凭证
// function getUploadToken() {
//   const options = {
//     scope: bucket,
//   };
//   const putPolicy = new qiniu.rs.PutPolicy(options);
//   return putPolicy.uploadToken(mac);
// }

// // 生成公开空间文件的访问链接
// function publicFileUrl(key) {
//   return `http://${domainOfBucket}/${key}`;
// }

// // 生成随机文件名
// function generateRandomFileName(fileName) {
//   const randomKey = uuidv4(); 
//   const parts = fileName.split(".");
//   const extension = parts.length > 1 ? `.${parts.pop()}` : '';
//   return `${randomKey}${extension}`; // 拼接随机字符串和文件扩展名
// }

// // 上传文件的函数
// export default async(data)=> {
//   const file = data.get("file")
//   if (!file) {
//     throw new Error("No file uploaded")
//   }
//   const newFilename = generateRandomFileName(file.name);
//   // 读取文件内容
//   // 将 File 对象转换为 Node.js 的 ReadableStream
//   const readableStream = Readable.fromWeb(file.stream());
//   // 读取文件内容到 Buffer
//   const chunks = [];
//   for await (const chunk of readableStream) {
//     chunks.push(chunk);
//   }
//   const buffer = Buffer.concat(chunks);
//   const uploadToken = getUploadToken(); // 获取上传凭证
//   const config = new qiniu.conf.Config();
//   const formUploader = new qiniu.form_up.FormUploader(config);

//   const putExtra = new qiniu.form_up.PutExtra(); // 额外参数
//   return new Promise((resolve, reject) => {
//     console.log(newFilename)
//     formUploader.put(uploadToken, newFilename, buffer, putExtra, (respErr, respBody, respInfo) => {
//       if (respErr) {
//         console.error('上传失败：', respErr);
//         reject(respErr);
//       } else {
//         if (respInfo.statusCode === 200) {
//           console.log('上传成功：', respBody);
//           const fileUrl = publicFileUrl(newFilename); // 生成公开空间文件的访问 URL
//           console.log('文件访问URL：', fileUrl);
//           resolve(fileUrl); // 返回文件的访问 URL
//           resolve(respBody);
//         } else {
//           console.error('上传失败：', respInfo);
//           reject(respInfo);
//         }
//       }
//     });
//   });
// }