import { useMutation } from '@tanstack/react-query';
import { docgenApi, type Document, type DocumentGenerateRequest } from '../api/docgen.api';

export const useGenerateDocument = () =>
  useMutation<Document, Error, DocumentGenerateRequest>({
    mutationFn: docgenApi.generate,
  });

export const TEMPLATES = [
  { value: 'indictment', label: 'Обвинительное заключение' },
  { value: 'summons', label: 'Повестка в трибунал' },
  { value: 'verdict', label: 'Приговор' },
] as const;

export type TemplateValue = (typeof TEMPLATES)[number]['value'];
