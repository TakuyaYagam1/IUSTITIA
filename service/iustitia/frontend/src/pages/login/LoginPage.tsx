import { DEFAULT_HOME, HOME_BY_ROLE } from '@/app/router/routes';
import { useAuthStore, userApi } from '@entities/user';
import { zodResolver } from '@hookform/resolvers/zod';
import { ApiError, type LoginRequest } from '@shared/api';
import { Button, Input, Paper } from '@shared/ui';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { z } from 'zod';
import styles from './LoginPage.module.css';

const schema = z.object({
  username: z.string().trim().min(1, 'Укажите логин'),
  password: z.string().min(1, 'Пароль обязателен'),
});

type FormValues = z.infer<typeof schema>;

export const LoginPage = (): JSX.Element => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const setSession = useAuthStore((s) => s.setSession);
  const [submitError, setSubmitError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { username: '', password: '' },
  });

  const onSubmit = async (values: FormValues) => {
    setSubmitError(null);
    try {
      const payload: LoginRequest = { username: values.username, password: values.password };
      const result = await userApi.login(payload);
      setSession({
        token: result.token,
        userId: result.user_id,
        role: result.role,
        dome: result.dome,
        username: values.username,
      });
      const next = searchParams.get('next');
      if (next) {
        window.location.href = next;
        return;
      }
      const landing = HOME_BY_ROLE[result.role] ?? DEFAULT_HOME;
      navigate(landing, { replace: true });
    } catch (err) {
      if (err instanceof ApiError) {
        setSubmitError(err.message);
      } else {
        setSubmitError('Связь с трибуналом прервана');
      }
    }
  };

  return (
    <div className={styles.wrap}>
      <div className={styles.card}>
        <div className={styles.brand}>
          <h1 className={styles.brandTitle}>IUSTITIA</h1>
          <span className={styles.brandSub}>Трибунал Одиннадцатого Государства</span>
        </div>

        <Paper variant="paper">
          <form className={styles.form} onSubmit={handleSubmit(onSubmit)} noValidate>
            <Input
              id="username"
              label="Удостоверение"
              autoComplete="username"
              placeholder="judge_3"
              error={errors.username?.message}
              {...register('username')}
            />
            <Input
              id="password"
              label="Пропускной код"
              type="password"
              autoComplete="current-password"
              error={errors.password?.message}
              {...register('password')}
            />

            {submitError && <div className={styles.errorBanner}>{submitError}</div>}

            <div className={styles.actions}>
              <Button type="submit" variant="primary" disabled={isSubmitting}>
                {isSubmitting ? 'Проверка...' : 'Войти'}
              </Button>
            </div>
          </form>
        </Paper>

        <div className={styles.footer}>Операция «Свободный Марс» · 2187</div>
      </div>
    </div>
  );
};
