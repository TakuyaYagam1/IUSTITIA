import { caseApi } from '@entities/case';
import { zodResolver } from '@hookform/resolvers/zod';
import { ApiError } from '@shared/api';
import { Button, Input, Textarea } from '@shared/ui';
import { useQuery } from '@tanstack/react-query';
import { useEffect, useMemo, useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { useComplaintSubmit } from '../model/useComplaintSubmit';
import styles from './ComplaintForm.module.css';

const schema = z.object({
  case_id: z.string().uuid('Выберите дело'),
  defendant: z.string().trim().min(2, 'Укажите обвиняемого'),
  article: z.string().trim().min(2, 'Укажите статью'),
  text: z.string().trim().min(10, 'Изложение не короче 10 символов'),
});

type FormValues = z.infer<typeof schema>;

interface ComplaintFormProps {
  onSuccess?: () => void;
}

export const ComplaintForm = ({ onSuccess }: ComplaintFormProps): JSX.Element => {
  const [serverError, setServerError] = useState<string | null>(null);

  const casesQuery = useQuery({
    queryKey: ['cases', 'search', 'all'],
    queryFn: () => caseApi.search({ limit: 100 }),
  });

  const mutation = useComplaintSubmit();

  const {
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { case_id: '', defendant: '', article: '', text: '' },
  });

  const selectedCaseId = watch('case_id');
  const cases = useMemo(() => casesQuery.data ?? [], [casesQuery.data]);

  useEffect(() => {
    if (!selectedCaseId) return;
    const selected = cases.find((c) => c.id === selectedCaseId);
    if (selected?.defendant) {
      setValue('defendant', selected.defendant, { shouldValidate: true, shouldDirty: true });
    }
  }, [selectedCaseId, cases, setValue]);

  const onSubmit = async (values: FormValues) => {
    setServerError(null);
    try {
      const body = `[${values.article}] против ${values.defendant}\n\n${values.text}`;
      await mutation.mutateAsync({ case_id: values.case_id, text: body });
      onSuccess?.();
    } catch (err) {
      if (err instanceof ApiError) {
        setServerError(err.message);
      } else {
        setServerError('Донос не принят. Повторите позже.');
      }
    }
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit(onSubmit)} noValidate>
      <div className={styles.selectField}>
        <label htmlFor="case_id" className={styles.label}>
          Дело
        </label>
        <select id="case_id" className={styles.select} {...register('case_id')}>
          <option value="">- выберите дело -</option>
          {cases.map((c) => (
            <option key={c.id} value={c.id}>
              №{c.seq_num}/Δ · {c.defendant}
            </option>
          ))}
        </select>
        {errors.case_id && <span className={styles.error}>{errors.case_id.message}</span>}
      </div>

      <div className={styles.grid}>
        <Input
          id="defendant"
          label="Обвиняемый"
          placeholder="Фамилия, имя, купол"
          error={errors.defendant?.message}
          {...register('defendant')}
        />
        <Input
          id="article"
          label="Статья"
          placeholder="УК 58-10"
          error={errors.article?.message}
          {...register('article')}
        />
      </div>

      <Textarea
        id="text"
        label="Изложение"
        placeholder="Краткое содержание обвинения, обстоятельства, свидетели..."
        rows={8}
        error={errors.text?.message}
        {...register('text')}
      />

      {serverError && <div className={styles.errorBanner}>{serverError}</div>}

      <div className={styles.actions}>
        <span className={styles.hint}>Каждый донос фиксируется в картотеке.</span>
        <Button type="submit" variant="primary" disabled={isSubmitting || mutation.isPending}>
          {mutation.isPending ? 'Передача...' : 'Подать донос'}
        </Button>
      </div>
    </form>
  );
};
