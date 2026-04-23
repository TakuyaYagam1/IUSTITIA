import { zodResolver } from '@hookform/resolvers/zod';
import { ApiError, type CaseOpinion, type PreliminaryVerdict } from '@shared/api';
import { Button, Textarea } from '@shared/ui';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { useOpinionFile } from '../model/useOpinionFile';
import styles from './OpinionForm.module.css';

const schema = z.object({
  preliminary_verdict: z.enum(['guilty', 'acquitted', 'dismissed']),
  reasoning: z.string().trim().min(10, 'Обоснование не короче 10 символов'),
});

type FormValues = z.infer<typeof schema>;

const VERDICT_LABELS: Record<PreliminaryVerdict, string> = {
  guilty: 'Виновен',
  acquitted: 'Невиновен',
  dismissed: 'Прекратить',
};

interface OpinionFormProps {
  caseId: string;
  onSuccess?: (o: CaseOpinion) => void;
}

export const OpinionForm = ({ caseId, onSuccess }: OpinionFormProps): JSX.Element => {
  const [serverError, setServerError] = useState<string | null>(null);
  const mutation = useOpinionFile(caseId);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { preliminary_verdict: 'guilty', reasoning: '' },
  });

  const onSubmit = async (values: FormValues) => {
    setServerError(null);
    try {
      const o = await mutation.mutateAsync(values);
      onSuccess?.(o);
    } catch (err) {
      if (err instanceof ApiError) {
        setServerError(err.message);
      } else {
        setServerError('Заключение не подано. Повторите позже.');
      }
    }
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit(onSubmit)} noValidate>
      <div className={styles.radioRow}>
        <span className={styles.label}>Предварительное заключение</span>
        <div className={styles.radios}>
          {(Object.keys(VERDICT_LABELS) as PreliminaryVerdict[]).map((v) => (
            <label key={v} className={styles.radio}>
              <input type="radio" value={v} {...register('preliminary_verdict')} />
              {VERDICT_LABELS[v]}
            </label>
          ))}
        </div>
        {errors.preliminary_verdict && (
          <span className={styles.error}>{errors.preliminary_verdict.message}</span>
        )}
      </div>

      <Textarea
        id="reasoning"
        label="Обоснование"
        placeholder="Обоснование заключения, улики, мотивы..."
        rows={8}
        error={errors.reasoning?.message}
        {...register('reasoning')}
      />

      {serverError && <div className={styles.errorBanner}>{serverError}</div>}

      <div className={styles.actions}>
        <span className={styles.hint}>Заключение передаётся судье.</span>
        <Button type="submit" variant="primary" disabled={isSubmitting || mutation.isPending}>
          {mutation.isPending ? 'Передача...' : 'Подать заключение'}
        </Button>
      </div>
    </form>
  );
};
