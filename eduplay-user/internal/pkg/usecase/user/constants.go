package user

const sendCodeTemplate = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Password Reset</title>
</head>
<body style="margin:0; padding:0; background-color:#f4f4f7; font-family:Arial, sans-serif;">

  <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#f4f4f7; padding:20px;">
    <tr>
      <td align="center">

        <table width="600" cellpadding="0" cellspacing="0" style="background-color:#ffffff; border-radius:8px; padding:30px;">
          
          <tr>
            <td align="center" style="padding-bottom:20px;">
              <h2 style="margin:0; color:#333;">EduPlay</h2>
            </td>
          </tr>

          <tr>
            <td style="color:#555; font-size:16px; line-height:1.5;">
              <p>Вы запросили сброс пароля.</p>
              <p>Используйте следующий код для подтверждения:</p>
            </td>
          </tr>

          <tr>
            <td align="center" style="padding:20px 0;">
              <div style="display:inline-block; padding:15px 25px; font-size:28px; letter-spacing:4px; background-color:#f4f4f7; border-radius:6px; font-weight:bold; color:#333;">
                {.Code}
              </div>
            </td>
          </tr>

          <tr>
            <td style="color:#555; font-size:14px;">
              <p>Код действителен в течение 10 минут.</p>
              <p>Если вы не запрашивали сброс пароля — просто проигнорируйте это письмо.</p>
            </td>
          </tr>

          <tr>
            <td style="padding-top:30px; font-size:12px; color:#999; text-align:center;">
              <p>© EduPlay</p>
            </td>
          </tr>

        </table>

      </td>
    </tr>
  </table>

</body>
</html>
`
