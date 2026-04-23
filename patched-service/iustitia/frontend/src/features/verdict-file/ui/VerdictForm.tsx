import { zodResolver } from '@hookform/resolvers/zod';
import { ApiError, type VerdictResult } from '@shared/api';
import { Button, Input, Textarea } from '@shared/ui';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { useVerdictFile } from '../model/useVerdictFile';
import styles from './VerdictForm.module.css';

type Verdict = 'guilty' | 'acquitted' | 'dismissed';

const schema = z
  .object({
    verdict: z.enum(['guilty', 'acquitted', 'dismissed']),
    sentence: z.string().trim().optional(),
    reasoning: z.string().trim().min(10, 'Обоснование не короче 10 символов'),
  })
  .refine((v) => v.verdict !== 'guilty' || (v.sentence && v.sentence.length >= 2), {
    message: 'Мера наказания обязательна при вердикте «виновен»',
    path: ['sentence'],
  });

type FormValues = z.infer<typeof schema>;

const VERDICT_LABELS: Record<Verdict, string> = {
  guilty: 'Виновен',
  acquitted: 'Оправдан',
  dismissed: 'Прекратить',
};

interface VerdictFormProps {
  caseId: string;
  onSuccess?: (res: VerdictResult) => void;
  disabled?: boolean;
  gateHint?: string;
}

export const VerdictForm = ({
  caseId,
  onSuccess,
  disabled,
  gateHint,
}: VerdictFormProps): JSX.Element => {
  const [serverError, setServerError] = useState<string | null>(null);
  const mutation = useVerdictFile(caseId);

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { verdict: 'guilty', sentence: '', reasoning: '' },
  });

  const currentVerdict = watch('verdict');
  const sentenceDisabled = currentVerdict !== 'guilty';

  const onSubmit = async (values: FormValues) => {
    setServerError(null);
    try {
      const payload = {
        verdict: values.verdict,
        reasoning: values.reasoning,
        sentence: values.verdict === 'guilty' ? values.sentence : undefined,
      };
      const res = await mutation.mutateAsync(payload);
      onSuccess?.(res);
    } catch (err) {
      if (err instanceof ApiError) {
        setServerError(err.message);
      } else {
        setServerError('Приговор не вынесен. Повторите позже.');
      }
    }
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit(onSubmit)} noValidate>
      <div className={styles.radioRow}>
        <span className={styles.label}>Приговор</span>
        <div className={styles.radios}>
          {(Object.keys(VERDICT_LABELS) as Verdict[]).map((v) => (
            <label key={v} className={styles.radio}>
              <input type="radio" value={v} {...register('verdict')} />
              {VERDICT_LABELS[v]}
            </label>
          ))}
        </div>
        {errors.verdict && <span className={styles.error}>{errors.verdict.message}</span>}
      </div>

      <Input
        id="sentence"
        label="Мера наказания"
        placeholder="Например: кислородный паёк 25 лет"
        disabled={sentenceDisabled}
        error={errors.sentence?.message}
        {...register('sentence')}
      />

      <Textarea
        id="reasoning"
        label="Обоснование суда"
        placeholder="Мотивировочная часть приговора..."
        rows={8}
        error={errors.reasoning?.message}
        {...register('reasoning')}
      />

      {serverError && <div className={styles.errorBanner}>{serverError}</div>}

      <div className={styles.actions}>
        <span className={styles.hint}>Приговор окончательный, обжалованию не подлежит.</span>
        <Button
          type="submit"
          variant="primary"
          disabled={isSubmitting || mutation.isPending || disabled}
        >
          {mutation.isPending ? 'Передача...' : 'Вынести приговор'}
        </Button>
      </div>
      {gateHint && <div className={styles.gateHint}>{gateHint}</div>}
    </form>
  );
};
