import { zodResolver } from '@hookform/resolvers/zod';
import { ApiError, type Case } from '@shared/api';
import { Button, Input, Textarea } from '@shared/ui';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { useCaseCreate } from '../model/useCaseCreate';
import styles from './CaseCreateForm.module.css';

const schema = z.object({
  defendant: z.string().trim().min(2, 'Укажите обвиняемого'),
  crime: z.string().trim().min(2, 'Укажите состав'),
  text: z.string().trim().min(10, 'Изложение не короче 10 символов'),
});

type FormValues = z.infer<typeof schema>;

interface CaseCreateFormProps {
  onSuccess?: (c: Case) => void;
}

export const CaseCreateForm = ({ onSuccess }: CaseCreateFormProps): JSX.Element => {
  const [serverError, setServerError] = useState<string | null>(null);
  const mutation = useCaseCreate();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { defendant: '', crime: '', text: '' },
  });

  const onSubmit = async (values: FormValues) => {
    setServerError(null);
    try {
      const c = await mutation.mutateAsync(values);
      onSuccess?.(c);
    } catch (err) {
      if (err instanceof ApiError) {
        setServerError(err.message);
      } else {
        setServerError('Дело не заведено. Повторите позже.');
      }
    }
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit(onSubmit)} noValidate>
      <div className={styles.grid}>
        <Input
          id="defendant"
          label="Обвиняемый"
          placeholder="Фамилия, имя, купол"
          error={errors.defendant?.message}
          {...register('defendant')}
        />
        <Input
          id="crime"
          label="Состав"
          placeholder="Краткая квалификация"
          error={errors.crime?.message}
          {...register('crime')}
        />
      </div>

      <Textarea
        id="text"
        label="Изложение"
        placeholder="Обстоятельства, свидетели, подозрения..."
        rows={8}
        error={errors.text?.message}
        {...register('text')}
      />

      {serverError && <div className={styles.errorBanner}>{serverError}</div>}

      <div className={styles.actions}>
        <span className={styles.hint}>Новое дело попадёт в канцелярию.</span>
        <Button type="submit" variant="primary" disabled={isSubmitting || mutation.isPending}>
          {mutation.isPending ? 'Передача...' : 'Завести дело'}
        </Button>
      </div>
    </form>
  );
};
